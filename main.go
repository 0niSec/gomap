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
				Name:     "ports",
				Aliases:  []string{"p"},
				Usage:    "Port ranges to scan (e.g. 80,443,8000-8100)",
				Category: "PORT SPECIFICATION:",
			},
			&cli.BoolFlag{
				Name:     "quiet",
				Aliases:  []string{"q"},
				Usage:    "Don't print the banner and other noise",
				Category: "OUTPUT MODES:",
			},
			&cli.StringFlag{
				Name:     "target",
				Aliases:  []string{"t"},
				Usage:    "The target to scan. Can accept IP addresses and domain names",
				Category: "TARGET SPECIFICATION:",
				Required: true,
			},
			&cli.DurationFlag{
				Name:     "timeout",
				Aliases:  []string{"T"},
				Usage:    "Timeout for the connection",
				Value:    10 * time.Second,
				Category: "TIMING AND PERFORMANCE:",
			},
			&cli.PathFlag{
				Name:     "output",
				Aliases:  []string{"o"},
				Usage:    "Output file",
				Category: "OUTPUT MODES:",
			},
			&cli.BoolFlag{
				Name:     "service",
				Aliases:  []string{"sV"},
				Usage:    "Probe open ports to determine service info",
				Value:    false,
				Category: "SERVICE/VERSION DETECTION:",
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
