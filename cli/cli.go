package main

import (
	"errors"
	"log"
	"os"

	"github.com/davecgh/go-spew/spew"
	"github.com/gbdubs/avibase_downloader"
	"github.com/urfave/cli/v2"
)

const (
	flag_region_code  = "region_code"
	flag_force_reload = "force_reload"
	flag_include_rare = "include_rare"
	flag_verbose      = "verbose"
)

func main() {
	app := &cli.App{
		Name:    "Avibase Downloader",
		Usage:   "A CLI for downloading lists of birds from the Avibase database.",
		Version: "1.0",
		Flags: []cli.Flag{
			&cli.StringSliceFlag{
				Name:  flag_region_code,
				Usage: "the region code (or codes) to download.",
			},
			&cli.BoolFlag{
				Name:  flag_include_rare,
				Usage: "If set, includes rare/incidental birds (as opposed to common birds).",
			},
			&cli.BoolFlag{
				Name:  flag_force_reload,
				Usage: "If set, will always look up from the server, rather than looking at a memoized version.",
			},
			&cli.BoolFlag{
				Name:  flag_verbose,
				Usage: "Whether to print the output or silently succeed, if the command succeeds.",
			},
		},
		Action: func(c *cli.Context) error {
			input := &avibase_downloader.Input{
				RegionCodes: c.StringSlice(flag_region_code),
				IncludeRare: c.Bool(flag_include_rare),
				ForceReload: c.Bool(flag_force_reload),
			}
			if len(input.RegionCodes) == 0 {
				return errors.New("one or more region_code must be provided")
			}
			output, err := input.Execute()
			if err != nil {
				return err
			}
			if c.Bool(flag_verbose) {
				spew.Dump(*output)
			}
			return nil
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
