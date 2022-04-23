package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

// https://api.slack.com/methods/conversations.history
func requestConversationHistory(token string, channelId string) (*http.Response, error) {
	url := "https://slack.com/api/conversations.history"

	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed http.NewRequest %w", err)
	}

	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	params := request.URL.Query()
	params.Add("channel", channelId)
	request.URL.RawQuery = params.Encode()

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("client.Do() error: %w", err)
	}

	return response, nil
}

func FetchSlackThreads(token string, channelId string) (SlackThreads, error) {
	response, err := requestConversationHistory(token, channelId)

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code is not 200: %s", response.Status)
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(response.Body)

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed ReadAll: %w", err)
	}

	var history ConversationHistoryJSON
	if err := json.Unmarshal(body, &history); err != nil {
		return nil, errors.New("can not unmarshal JSON")
	}

	if !history.Ok {
		return nil, fmt.Errorf("なんかしっぱい: error: %v ResponseMetaData: %v", history.Error, history.ResponseMetadata)
	}

	return history.ToSlackThreads(), nil
}
