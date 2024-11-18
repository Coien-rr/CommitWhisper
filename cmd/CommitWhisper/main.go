package main

import (
	"fmt"
	"log"

	"github.com/Coien-rr/CommitWhisper/internal/git"
	"github.com/Coien-rr/CommitWhisper/internal/models"
)

func main() {
	diff, err := git.GetGitDiff()
	if err != nil {
		log.Fatal(err)
	}
	res, err := models.GetModelResponse(
		"https://dashscope.aliyuncs.com/compatible-mode/v1/chat/completions",
		"qwen2.5-coder-3b-instruct",
		"sk-502bbd61da624094920ce1c00375f45f",
		diff,
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Hello, CommitWhisper")
	fmt.Print(res)
}
