package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Coien-rr/CommitWhisper/internal/git"
	"github.com/Coien-rr/CommitWhisper/internal/whisper"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "commitwhisper",
		Usage: "Generate AI Commit By LLM Using Git Diff Info",
		Action: func(*cli.Context) error {
			config := whisper.GetConfig()

			diff, err := git.GetGitDiff()
			if err != nil {
				log.Fatal(err)
			}

			if w := whisper.NewWhisper(config); w != nil {
				w.Greet()
				w.HandleGeneratedCommitMsg(diff)
			} else {
				return nil
			}
			return nil
		},
		Commands: []*cli.Command{
			{
				Name:    "config",
				Aliases: []string{"c"},
				Usage:   "Reconfig LLM Config, Such as AiProvider, Model, Url, Key",
				Action: func(cCtx *cli.Context) error {
					fmt.Println("completed task: ", cCtx.Args().First())
					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
