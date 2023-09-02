package main

import (
	"fmt"
	"log"
	"os"

	"github.com/mjc-gh/fickle_dns/server/api"
	"github.com/urfave/cli/v2"
)

func main() {
	var port int

	app := &cli.App{
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:        "port",
				Value:       3036,
				Usage:       "port number",
				Destination: &port,
			},
		},
		Action: func(cCtx *cli.Context) error {
			r := api.Setup()
			r.Run(fmt.Sprintf("127.0.0.1:%d", port))

			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
