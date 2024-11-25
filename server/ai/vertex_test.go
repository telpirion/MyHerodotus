package ai

import (
	"strings"
	"testing"
)

func TestCreatePrompt(t *testing.T) {
	query := "I'm a query"
	got, err := createPrompt(query, GeminiTemplate)
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(got, query) {
		t.Errorf("got: %v, want: %v", got, query)
	}
}

func TestSetConversationContext(t *testing.T) {
	convoHistory := []ConversationBit{
		{
			UserQuery:   "test user query",
			BotResponse: "test bot response",
		},
		{
			UserQuery:   "test user query 2",
			BotResponse: "test bot response 2",
		},
	}
	err := SetConversationContext(convoHistory)
	if err != nil {
		t.Fatal(err)
	}

	got := cachedContext

	if !strings.Contains(got, convoHistory[0].UserQuery) {
		t.Errorf("got: %v, want: %v", cachedContext, convoHistory[0].UserQuery)
	}
}
