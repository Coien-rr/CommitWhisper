package main

import (
	"log"
	"os"

	"github.com/Coien-rr/CommitWhisper/internal/whisper"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "cw",
		Usage: "Generate AI Commit By LLM Using Git Diff Info",
		Action: func(*cli.Context) error {
			config := whisper.GetConfig()

			if w := whisper.NewWhisper(config); w != nil {
				w.Run()
			} else {
				return nil
			}
			return nil
		},
		Commands: []*cli.Command{
			{
				Name:    "reconfig",
				Aliases: []string{"rc"},
				Usage:   "Reconfig LLM Config, Such as AiProvider, Model, Url, Key",
				Action: func(cCtx *cli.Context) error {
					whisper.ReConfig()
					return nil
				},
			},
			{
				Name:    "showconfig",
				Aliases: []string{"sc"},
				Usage:   "Show LLM Config",
				Action: func(cCtx *cli.Context) error {
					whisper.ShowConfig()
					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
