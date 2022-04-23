package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// https://api.slack.com/methods/conversations.replies
func requestConversationReplies(token, channelId, ts string) (*http.Response, error) {
	url := "https://slack.com/api/conversations.replies"

	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed http.NewRequest %w", err)
	}
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	params := request.URL.Query()
	params.Add("channel", channelId)
	params.Add("ts", ts)
	request.URL.RawQuery = params.Encode()

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("client.Do() error: %w", err)
	}

	return response, nil
}

func FetchSlackFiles(token, channelId, ts string) (SlackFiles, error) {
	response, err := requestConversationReplies(token, channelId, ts)

	_, _ = response, err

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code is not 200: %s", response.Status)
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed ReadAll: %w", err)
	}

	var replies ConversationRepliesJSON
	if err := json.Unmarshal(body, &replies); err != nil {
		return nil, fmt.Errorf("can not unmarshal JSON: %w", err)
	}

	return replies.ToSlackFiles(), nil
}
