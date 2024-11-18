package main

import (
	"fmt"
	"os"
	"strings"
	"testing"
)

func TestDeidentify(t *testing.T) {
	projectID := os.Getenv("PROJECT_ID")
	firstName := "Artemis"
	lastName := "Schmidt"
	email := "xanthippe@example.com"
	testInput := fmt.Sprintf(`
My name is %s %s and my email is %s. I want to go see the Great Pyramids of Egypt.
`, firstName, lastName, email)

	got, err := deidentify(projectID, testInput)
	if err != nil {
		t.Fatal(err)
	}

	if strings.Contains(strings.ToLower(got), firstName) ||
		strings.Contains(strings.ToLower(got), lastName) ||
		strings.Contains(strings.ToLower(got), email) {
		t.Errorf("No deidentification; got %s", got)
	}
}
