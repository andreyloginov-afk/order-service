package builder

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"reflect"
	"sync"
	"syscall"

	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"

	"github.com/andreyloginov-afk/order-service/internal/app/config"
	rhandler "github.com/andreyloginov-afk/order-service/internal/app/handler/http"
	rhealth "github.com/andreyloginov-afk/order-service/internal/app/handler/http/health"
	horder "github.com/andreyloginov-afk/order-service/internal/app/handler/http/order"
	"github.com/andreyloginov-afk/order-service/internal/app/processor"
	rprocessor "github.com/andreyloginov-afk/order-service/internal/app/processor/http"
	"github.com/andreyloginov-afk/order-service/internal/app/repository"
	rcpostgres "github.com/andreyloginov-afk/order-service/internal/app/repository/conn/postgres"
	porder "github.com/andreyloginov-afk/order-service/internal/app/repository/order"
	"github.com/andreyloginov-afk/order-service/internal/app/service"
	sorder "github.com/andreyloginov-afk/order-service/internal/app/service/order"
)

type Builder struct {
	cCtx *cli.Context
	ctx  context.Context
	wg   sync.WaitGroup
	err  error
	cfg  config.Config

	chErrors chan error

	connPostgres *rcpostgres.Client

	orderRepo    repository.Order
	orderService service.Order

	healthHandler rhandler.Health
	orderHandler  rhandler.Order

	processors []processor.Processor
}

func NewBuilder(cCtx *cli.Context) *Builder {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())

	b := &Builder{
		cCtx:          cCtx,
		ctx:           ctx,
		chErrors:      make(chan error, 16),
		healthHandler: rhealth.NewHandler(),
	}

	go b.waitForSignal(ch, cancel)

	return b
}

func (b *Builder) BuildConfig() *Builder {
	return b.exec(func(b *Builder) {
		config.Load(config.LoadArgs{
			Output:          os.Stdout,
			EnableSimpleLog: b.cCtx.Bool("no-json"),
		})
		b.cfg = config.Root
	})
}

func (b *Builder) Run() {
	if b.ctx.Err() != nil {
		log.Info().Msg("Shutdown during initialization")
		return
	}

	if b.err != nil {
		log.Fatal().Err(b.err).Msg("Failed to initialize application")
	}

	if b.connPostgres != nil {
		processor.WatchForShutdown(b.ctx, &b.wg, b.connPostgres)
	}

	log.Info().Msg("Application initialized")
	defer log.Info().Msg("Application completed")

	for _, p := range b.processors {
		p.StartAsync(b.ctx, &b.wg)
	}

	b.wg.Wait()
}

////////////////////////////////////////////////////////////////////////////////
///// REPOSITORY CONNECTIONS ///////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

func (b *Builder) BuildRepoConnPostgres() *Builder {
	return b.exec(func(b *Builder) {
		client, err := rcpostgres.NewClient(b.ctx, b.cfg.Repository.Postgres)
		if err != nil {
			b.err = fmt.Errorf("postgres: %w", err)
			return
		}
		b.connPostgres = client
	})
}

////////////////////////////////////////////////////////////////////////////////
///// REPOSITORIES /////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

func (b *Builder) BuildRepoOrder() *Builder {
	return b.exec(func(b *Builder) {
		b.orderRepo = porder.NewRepo(b.connPostgres)
	}, b.connPostgres)
}

////////////////////////////////////////////////////////////////////////////////
///// SERVICES /////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

func (b *Builder) BuildServiceOrder() *Builder {
	return b.exec(func(b *Builder) {
		b.orderService = sorder.NewService(b.orderRepo)
	}, b.orderRepo)
}

////////////////////////////////////////////////////////////////////////////////
///// HANDLERS /////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

func (b *Builder) BuildHandlerOrder() *Builder {
	return b.exec(func(b *Builder) {
		b.orderHandler = horder.NewHandler(b.orderService)
	}, b.orderService)
}

////////////////////////////////////////////////////////////////////////////////
///// PROCESSORS ///////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

func (b *Builder) BuildProcHttp() *Builder {
	return b.exec(func(b *Builder) {
		p := rprocessor.NewHTTP(b.healthHandler, b.orderHandler, b.cfg.Processor.WebServer)
		b.processors = append(b.processors, p)
	}, b.orderHandler)
}

////////////////////////////////////////////////////////////////////////////////
///// PRIVATE //////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

func (b *Builder) waitForSignal(ch chan os.Signal, cancel context.CancelFunc) {
	defer signal.Stop(ch)
	select {
	case sig := <-ch:
		log.Info().Str("signal", sig.String()).Msg("Shutdown is requested")
		cancel()
	case <-b.ctx.Done():
	}
}

func (b *Builder) exec(fn func(*Builder), deps ...any) *Builder {
	if b.err != nil {
		return b
	}
	for _, dep := range deps {
		if isNilDep(dep) {
			b.err = errors.New("required dependency is nil")
			return b
		}
	}
	fn(b)
	return b
}

func isNilDep(v any) bool {
	if v == nil {
		return true
	}
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Ptr, reflect.Interface, reflect.Map, reflect.Slice, reflect.Chan, reflect.Func:
		return rv.IsNil()
	}
	return false
}
