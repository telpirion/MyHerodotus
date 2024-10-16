package main

import (
	"context"
	"log"

	"cloud.google.com/go/vertexai/genai"
)

// textPredict generates text with certain prompt and configurations.
func textPredict(message, projectID, model string) (string, error) {
	ctx := context.Background()
	location := "us-west1"

	client, err := genai.NewClient(ctx, projectID, location)
	if err != nil {
		log.Println(err)
		return "", err
	}
	defer client.Close()

	llm := client.GenerativeModel(model)
	prompt := genai.Text(
		`The user wants to go on a vacation somewhere. They want your opinion on what to do. Here is their question: ` + message,
	)

	resp, err := llm.GenerateContent(ctx, prompt)
	if err != nil {
		log.Println(err)
		return "", err
	}

	candidate := resp.Candidates[0].Content.Parts[0].(genai.Text)
	return string(candidate), nil
}
