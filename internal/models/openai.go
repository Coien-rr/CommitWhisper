package models

import (
	"encoding/json"
	"fmt"

	"github.com/Coien-rr/CommitWhisper/internal/comm"
)

type OpenAIModel struct {
	BaseModel
	localSession  session
	isRefineStage bool
}

func (m *OpenAIModel) addPrompt(promptMsg string) {
	if m.isRefineStage {
		m.localSession.appendMessage("user", getCommitGeneratePrompt(promptMsg))
		m.isRefineStage = true
	} else {
		m.localSession.appendMessage("user", getCommitRefinePrompt(promptMsg))
	}
}

func (m *OpenAIModel) addModelResponse(respMsg string) {
	m.localSession.appendMessage("assistant", respMsg)
}

func (m *OpenAIModel) prepareSessionChatReqBody() genericLLMsServiceReqBody {
	return genericLLMsServiceReqBody{
		Model:    m.modelName,
		Messages: m.localSession.getMessages(),
	}
}

func NewOpenAIModelAgent(modelName, baseUrl, apiKey string) (Model, error) {
	model := &OpenAIModel{
		BaseModel:     BaseModel{modelName: modelName, url: baseUrl, key: apiKey},
		isRefineStage: false,
	}

	err := model.initSession()

	return model, err
}

// NOTE: Local Session
func (m *OpenAIModel) initSession() error {
	m.localSession.appendMessage("system", GetSystemPrompt())
	return nil
}

func (m *OpenAIModel) GenerateCommitMessage(diffInfo string) (string, error) {
	m.addPrompt(diffInfo)

	client := comm.NewLLMsServiceClient(m.key, m.url)
	requestBody, err := json.Marshal(m.prepareSessionChatReqBody())
	if err != nil {
		return "", fmt.Errorf(
			"ERROR(GenerateCommitMessage): Failed to marshal request body: %w",
			err,
		)
	}

	resp, statusCode, err := client.CreateChatReqWithLLMs(requestBody)
	if err != nil {
		return "", fmt.Errorf(
			"ERROR(GenerateCommitMessage): %w",
			err,
		)
	}

	msg, err := handleChatRespFromLLMs(resp, statusCode)
	if err != nil {
		return msg, err
	}

	m.addModelResponse(msg)

	return msg, nil
}
