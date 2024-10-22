package main

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
