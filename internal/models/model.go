package models

import (
	"fmt"
	"net/http"
)

type Model interface {
	PrepareRequest(diffInfo string) (*http.Request, error)
	CreateContextSession() (string, error)
	CreateSessionChatRequest(diffInfo string) (*http.Request, error)
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
		return &DoubaoModel{
			BaseModel: BaseModel{modelName: modelName, url: url, key: key},
			ContextID: "",
			chatCount: 0,
		}, nil
	default:
		return nil, fmt.Errorf(
			"CreateModelError: %v is unsupported yet, Coming Soon  ",
			aiProvider,
		)
	}
}
