package main

import (
	"github.com/urfave/cli/v2"
	"html-text-extractor/internal"
	"log"
	"os"
	"time"
)

func main() {
	app := &cli.App{
		Name:  "HTML text extractor",
		Usage: "Use to extract all text from HTML page",
		Action: func(context *cli.Context) error {

			ignorelist := context.StringSlice("ignore")

			ignoreMap := make(map[string]struct{})

			for _, item := range ignorelist {
				ignoreMap[item] = struct{}{}
			}

			err := internal.ScrapeDomainsFromFile(
				context.String("input"),
				context.Duration("timeout"),
				context.Int("concurrency"),
				ignoreMap,
				context.String("output"),
			)
			return err
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "input",
				Value: "domains.txt",
				Usage: "File with list of domains to scrape seperated by new line",
			},
			&cli.DurationFlag{
				Name:  "timeout",
				Value: time.Second * 10,
				Usage: "Request timeout",
			},
			&cli.IntFlag{Name: "concurrency", Value: 10, Usage: "Max concurrent http requests"},
			&cli.StringSliceFlag{Name: "ignore", Usage: "HTML tags to ignore"},
			&cli.StringFlag{Name: "output", Value: "out.csv", Usage: "Output CSV file name"},
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
