package comm

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

type client struct {
	httpClient *http.Client
	apiKey     string
	baseUrl    string
}

var (
	clientInstance *client
	once           sync.Once
)

func NewLLMsServiceClient(apiKey, baseUrl string) *client {
	once.Do(func() {
		clientInstance = &client{
			httpClient: &http.Client{
				Timeout: 30 * time.Second,
			},
			apiKey:  apiKey,
			baseUrl: baseUrl,
		}
	})
	return clientInstance
}

func (c *client) CreateLLMsContextSession(reqBody []byte) ([]byte, int, error) {
	req, err := c.CreateNewSessionRequest(reqBody)
	if err != nil {
		return nil, 0, fmt.Errorf(
			"ERROR(CreateLLMsContextSession): %w",
			err,
		)
	}

	return c.fetchLLMsServiceResp(req)
}

func (c *client) CreateChatReqWithLLMs(reqBody []byte) ([]byte, int, error) {
	req, err := c.CreateSessionChatRequest(reqBody)
	if err != nil {
		return nil, 0, fmt.Errorf(
			"ERROR(CreateChatReqWithLLMs): %w",
			err,
		)
	}

	return c.fetchLLMsServiceResp(req)
}
