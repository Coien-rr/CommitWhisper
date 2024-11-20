package models

import "net/http"

type Model interface {
	PrepareRequest(diffInfo string) (*http.Request, error)
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
	url       string
	modelName string
	key       string
}

func prepareQuestionContent(diffInfo string) string {
	return "Please write a commit message for these git changes, " + diffInfo
}

func CreateModel(llmModel, url, modelName, key string) Model {
	switch llmModel {
	case "qwen":
		return &QWENModel{BaseModel{url: url, modelName: modelName, key: key}}
	}
	return nil
}
