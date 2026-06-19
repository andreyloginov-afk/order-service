package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"

	"github.com/andreyloginov-afk/order-service/cmd"
)

func main() {
	app := &cli.App{
		Name:    "order-service",
		Version: "1.0.0",
		Usage:   "Order management service",
		Commands: []*cli.Command{
			cmd.WebServer(),
		},
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "no-json",
				Usage: "disable JSON logging, use plain text output",
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}
