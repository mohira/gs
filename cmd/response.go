package cmd

import (
	"strings"
)

type SlackThreads []SlackThread

type SlackThread struct {
	Text string
	Ts   string
}

func (t *SlackThread) FirstLine() string {
	return strings.Split(t.Text, "\n")[0]
}

type ConversationHistoryJSON struct {
	Ok               bool   `json:"ok"`
	Error            string `json:"error"`
	ResponseMetadata struct {
		Messages []string `json:"messages"`
	} `json:"response_metadata"`
	Messages []Message `json:"messages"`
}

func (h *ConversationHistoryJSON) ToSlackThreads() SlackThreads {
	var slackThreads []SlackThread

	for _, m := range h.Messages {
		if m.IsTextMessage() && m.IsThread() {
			t := SlackThread{
				Text: m.Text,
				Ts:   m.ThreadTs,
			}
			slackThreads = append(slackThreads, t)
		}
	}

	return slackThreads
}

type Message struct {
	Type     string `json:"type"`
	Text     string `json:"text"`
	User     string `json:"user"`
	ThreadTs string `json:"thread_ts,omitempty"`
}

func (m *Message) IsThread() bool {
	return m.ThreadTs != ""
}

func (m *Message) IsTextMessage() bool {
	return m.Type == "message"
}
