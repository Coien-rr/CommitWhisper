package main

import (
	"log"

	"github.com/Coien-rr/CommitWhisper/internal/git"
	"github.com/Coien-rr/CommitWhisper/internal/interaction"
	"github.com/Coien-rr/CommitWhisper/internal/models"
)

func main() {
	interaction.Greet()

	diff, err := git.GetGitDiff()
	if err != nil {
		log.Fatal(err)
	}

	model := models.CreateModel(
		"https://dashscope.aliyuncs.com/compatible-mode/v1/chat/completions",
		"qwen2.5-coder-3b-instruct",
		`sk-502bbd61da624094920ce1c00375f45f`,
	)

	model.HandleGeneratedCommitMsg(diff)
}
