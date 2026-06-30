package cmd

import (
	"strings"

	"github.com/urfave/cli/v2"

	"github.com/andreyloginov-afk/order-service/internal/app/builder"
)

func WebServer() *cli.Command {
	return &cli.Command{
		Name:    "web-server",
		Aliases: []string{"ws"},
		Usage:   "Start HTTP server with all routes",
		Description: strings.TrimSpace(`
Initializes all dependencies (config, DB, repositories, services, handlers)
and starts the HTTP server. Graceful shutdown on SIGINT/SIGTERM.
`),
		Action:          cmdWebServer,
		HideHelpCommand: true,
	}
}

func cmdWebServer(cCtx *cli.Context) error {
	builder.NewBuilder(cCtx).
		BuildConfig().
		BuildRepoConnPostgres().
		BuildRepoOrder().
		BuildServiceOrder().
		BuildHandlerOrder().
		BuildProcHttp().
		Run()
	return nil
}
