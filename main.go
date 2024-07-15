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
		Name:  "gomap",
		Usage: "The Go port scanner",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "ports",
				Aliases: []string{"p"},
				Usage:   "Port ranges to scan (e.g. 80,443,8000-8100)",
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
		Before: func(c *cli.Context) error {
			if !c.Bool("quiet") {
				PrintBanner()
				ScanInfo(c)
			}
			return nil
		},
		Action:               gomapcli.Runner,
		EnableBashCompletion: true,
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}

}
