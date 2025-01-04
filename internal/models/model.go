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

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type BaseModel struct {
	modelName string
	url       string
	key       string
}

func CreateModel(aiProvider, modelName, url, key string) (Model, error) {
	switch aiProvider {
	// case "Qwen":
	// 	return &QWENModel{BaseModel{modelName: modelName, url: url, key: key}}, nil
	// case "OpenAI":
	// 	return &OpenAIModel{BaseModel{modelName: modelName, url: url, key: key}}, nil
	case "Doubao":
		modelAgent, err := NewDoubaoModel(modelName, url, key)
		if err != nil {
			return nil, fmt.Errorf(
				"CreateModelError: Create %v Model Failed For %w",
				aiProvider, err,
			)
		}
		return modelAgent, nil
	default:
		return nil, fmt.Errorf(
			"CreateModelError: %v is unsupported yet, Coming Soon  ",
			aiProvider,
		)
	}
}
