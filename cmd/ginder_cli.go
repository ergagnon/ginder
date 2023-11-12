package main

import (
	"os"

	"github.com/ergagnon/ginder/cmd/internal"
	"github.com/ergagnon/gocmder"
)

func main2() {
	cli, err := gocmder.NewCmder(internal.AppConfig{}, func(cfg any) {
			app := internal.NewApp(cfg.(internal.AppConfig))
			app.Run()
		},
		gocmder.WithPrefix("GINDER"),
		gocmder.WithLongDesc("\nGINDER help you find your most sensitive informations.\n"),
		gocmder.WithVersion("0.0.1"),
	)

	if err != nil {
		os.Exit(1)
	}

	if err := cli.Execute(); err != nil {
		os.Exit(1)
	}
}