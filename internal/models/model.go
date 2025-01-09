package models

import (
	"fmt"
)

type Model interface {
	GenerateCommitMessage(diffInfo string) (string, error)
}

type RequestBody struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type BaseModel struct {
	modelName string
	url       string
	key       string
}

func CreateModel(aiProvider, modelName, url, key string) (Model, error) {
	var modelAgent Model
	var err error
	switch aiProvider {
	case "Qwen":
		modelAgent, err = NewQwenModelAgent(modelName, url, key)
	case "Doubao":
		modelAgent, err = NewDoubaoModelAgent(modelName, url, key)
	// TODO: OpenAI
	case "OpenAI":
		modelAgent, err = NewDoubaoModelAgent(modelName, url, key)
	default:
		return nil, fmt.Errorf(
			"CreateModelError: %v is unsupported yet, Coming Soon  ",
			aiProvider,
		)
	}

	if err != nil {
		return nil, fmt.Errorf(
			"CreateModelError: Create %v Model Failed For %w",
			aiProvider, err,
		)
	}
	return modelAgent, err
}
