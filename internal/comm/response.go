package comm

import (
	"fmt"
	"io"
	"net/http"
)

func (c *client) fetchLLMsServiceResp(req *http.Request) ([]byte, int, error) {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf(
			"ERROR(CreateLLMsContextSession): %w",
			err,
		)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, fmt.Errorf(
			"ERROR(handleSessionCreateResponse): Failed to read response body: %w",
			err,
		)
	}

	return body, resp.StatusCode, nil
}
