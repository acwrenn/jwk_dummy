package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v2"

	"github.com/acwrenn/jwk_dummy/internal/server"
)

func main() {
	app := &cli.App{
		Name:   "Dummy JWK Server - I'll sign anything!",
		Usage:  "Run a little JWK server, that signs any payload and serves it's keys back to validate the sig.",
		Action: runServer,
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:  "port",
				Value: 3333,
				Usage: "Port to listen on.",
			},
			&cli.StringFlag{
				Name:  "address",
				Value: "localhost",
				Usage: "Address to listen on.",
			},
			&cli.StringFlag{
				Name:  "config-route",
				Value: "/.well-known/openid-configuration",
				Usage: "Can be used to set the initial configuation URL.",
			},
			&cli.StringFlag{
				Name:  "key-file",
				Value: "",
				Usage: "Can be used to load a pre-used set of keys from a file.",
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func runServer(c *cli.Context) error {
	address := fmt.Sprintf("%s:%d", c.String("address"), c.Int("port"))
	config := server.Config{
		ConfigRoute: c.String("config-route"),
		KeyFile:     c.String("key-file"),

		Address:  c.String("address"),
		Port:     c.Int("port"),
		Protocol: "http",
	}
	return server.Run(address, config)
}
