package builder

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"

	"github.com/andreyloginov-afk/order-service/internal/app/config"
	rhandler "github.com/andreyloginov-afk/order-service/internal/app/handler/http"
	rhealth "github.com/andreyloginov-afk/order-service/internal/app/handler/http/health"
	"github.com/andreyloginov-afk/order-service/internal/app/processor"
	rprocessor "github.com/andreyloginov-afk/order-service/internal/app/processor/http"
	rcpostgres "github.com/andreyloginov-afk/order-service/internal/app/repository/conn/postgres"
)

type Builder struct {
	cCtx *cli.Context
	ctx  context.Context
	wg   sync.WaitGroup
	err  error
	cfg  config.Config

	chErrors chan error

	connPostgres *rcpostgres.Client

	healthHandler rhandler.Health

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
	return b.exec(b.buildConfig)
}

func (b *Builder) Run() error {
	if b.err != nil {
		return nil
	}

	if b.connPostgres != nil {
		processor.WatchForShutdown(b.ctx, &b.wg, b.connPostgres)
	}

	log.Info().Msg("Application initialized")

	for _, p := range b.processors {
		p.StartAsync(b.ctx, &b.wg)
	}

	b.wg.Wait()
	b.printErrors()
	log.Info().Msg("Application completed")
	return b.err
}

////////////////////////////////////////////////////////////////////////////////
///// REPOSITORY CONNECTIONS ///////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

func (b *Builder) BuildRepoConnPostgres() *Builder {
	return b.exec(func() error {
		client, err := rcpostgres.NewClient(b.ctx, b.cfg.Repository.Postgres)
		if err != nil {
			return fmt.Errorf("postgres: %w", err)
		}
		b.connPostgres = client
		return nil
	})
}

////////////////////////////////////////////////////////////////////////////////
///// PROCESSORS ///////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

func (b *Builder) BuildProcHttp() *Builder {
	return b.exec(func() error {
		p := rprocessor.NewHTTP(b.healthHandler, b.cfg.Processor.WebServer)
		b.processors = append(b.processors, p)
		return nil
	})
}

////////////////////////////////////////////////////////////////////////////////
///// PRIVATE //////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

func (b *Builder) buildConfig() error {
	config.Load(config.LoadArgs{
		Output:          os.Stdout,
		EnableSimpleLog: b.cCtx.Bool("no-json"),
	})
	b.cfg = config.Root
	return nil
}

func (b *Builder) waitForSignal(ch chan os.Signal, cancel context.CancelFunc) {
	defer signal.Stop(ch)
	select {
	case sig := <-ch:
		log.Info().Str("signal", sig.String()).Msg("Shutdown is requested")
		cancel()
	case <-b.ctx.Done():
	}
}

func (b *Builder) printErrors() {
	close(b.chErrors)
	for err := range b.chErrors {
		log.Error().Err(err).Send()
	}
}

func (b *Builder) exec(fn func() error) *Builder {
	if b.err != nil {
		return b
	}
	b.err = fn()
	return b
}
