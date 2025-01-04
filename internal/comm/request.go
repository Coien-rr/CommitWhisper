package comm

import (
	"bytes"
	"fmt"
	"net/http"
)

func CreateLLMsRequest(reqBody []byte, key, url string) (*http.Request, error) {
	req, err := http.NewRequest(
		http.MethodPost,
		url,
		bytes.NewBuffer(reqBody),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+key)
	req.Header.Set("Content-Type", "application/json")

	return req, nil
}

func (c *client) createLLMsRequest(reqBody []byte, endPointPath string) (*http.Request, error) {
	req, err := http.NewRequest(
		http.MethodPost,
		c.baseUrl+endPointPath,
		bytes.NewBuffer(reqBody),
	)
	if err != nil {
		return nil, fmt.Errorf("ERROR(createLLMsRequest): failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	return req, nil
}

func (c *client) CreateNewSessionRequest(createSessionReqBody []byte) (*http.Request, error) {
	return c.createLLMsRequest(createSessionReqBody, "/context/create")
}

func (c *client) CreateSessionChatRequest(sessionChatReqBody []byte) (*http.Request, error) {
	return c.createLLMsRequest(sessionChatReqBody, "/context/chat/completions")
}
