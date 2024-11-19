package models

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

func GetModelResponse(url, model, key, diffInfo string) (string, error) {
	reqBody := PrepareRequestBody(model, diffInfo)

	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to encode request body: %v", err)
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(reqBytes))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+key)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("received non-200 response: %v", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %v", err)
	}

	var response ResponseBody
	if err := json.Unmarshal(body, &response); err != nil {
		return "", fmt.Errorf("failed to parse response JSON: %v", err)
	}

	return response.Choices[0].Message.Content, nil
}
