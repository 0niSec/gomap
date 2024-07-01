package main

import (
	"log"
	"os"
	"time"

	"github.com/0niSec/gomap/gomapcli"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:      "gomap",
		Usage:     "The lightweight, thin Go port scanning tool",
		Version:   "0.1.0",
		Copyright: "(c) 2024 0niSec. All rights reserved.",
		Authors: []*cli.Author{
			&cli.Author{
				Name: "0niSec",
			},
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "ports",
				Aliases: []string{"p"},
				Usage:   "Port ranges to scan",
			},
			&cli.BoolFlag{
				Name:    "quiet",
				Aliases: []string{"q"},
				Usage:   "Don't print the banner and other noise",
			},
			&cli.StringFlag{
				Name:    "target",
				Aliases: []string{"t"},
				Usage:   "The target to scan",
			},
			&cli.DurationFlag{
				Name:    "timeout",
				Aliases: []string{"T"},
				Usage:   "Timeout for the connection",
				Value:   10 * time.Second,
			},
			&cli.PathFlag{
				Name:    "output",
				Aliases: []string{"o"},
				Usage:   "Output file",
			},
		},
		Action: gomapcli.Runner,
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
