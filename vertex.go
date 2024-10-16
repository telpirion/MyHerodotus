package main

import (
	"context"
	"fmt"
	"log"
	"os"

	aiplatform "cloud.google.com/go/aiplatform/apiv1"
	"cloud.google.com/go/aiplatform/apiv1/aiplatformpb"
	"cloud.google.com/go/vertexai/genai"
	"google.golang.org/api/option"
	"google.golang.org/protobuf/types/known/structpb"
)

func createPrompt(message string) string {
	return `
You are a helpful travel agent. The user wants to have a conversation with you about where to go.
You are going to help them plan their trip.

Example:
[USER]: I want to go to Italy.
[RESPONSE]: Italy is a fantastic place to go! What would you like to experience: the food, the
history, the art, the people, or some combination?

Here is the user query. Respond to the user's request. Check your answer before responding.

[USER]:` + message
}

// textPredictGemma2 generates text using a Gemma2 hosted model
func textPredictGemma(message, projectID string) (string, error) {
	ctx := context.Background()
	location := "us-west1"
	endpointID := os.Getenv("ENDPOINT_ID")
	gemma2Endpoint := fmt.Sprintf("projects/%s/locations/%s/endpoints/%s", projectID, location, endpointID)

	apiEndpoint := fmt.Sprintf("%s-aiplatform.googleapis.com:443", location)
	client, err := aiplatform.NewPredictionClient(ctx, option.WithEndpoint(apiEndpoint))
	if err != nil {
		return "", fmt.Errorf("unable to create prediction client: %v", err)
	}
	defer client.Close()

	parameters := map[string]interface{}{
		"temperature":     0.5,
		"maxOutputTokens": 1024,
		"topP":            0.2,
		"topK":            10,
	}

	promptValue, err := structpb.NewValue(map[string]interface{}{
		"inputs":     createPrompt(message),
		"parameters": parameters,
	})
	if err != nil {
		return "", err
	}

	req := &aiplatformpb.PredictRequest{
		Endpoint:  gemma2Endpoint,
		Instances: []*structpb.Value{promptValue},
	}

	resp, err := client.Predict(ctx, req)
	if err != nil {
		return "", err
	}

	prediction := resp.GetPredictions()
	value := prediction[0].GetStringValue()

	return value, nil
}

// textPredictGemini generates text using a Gemini 1.5 Flash model
func textPredictGemini(message, projectID string) (string, error) {
	ctx := context.Background()
	location := "us-west1"
	model := "gemini-1.5-flash-001"

	client, err := genai.NewClient(ctx, projectID, location)
	if err != nil {
		log.Println(err)
		return "", err
	}
	defer client.Close()

	llm := client.GenerativeModel(model)
	prompt := createPrompt(message)

	resp, err := llm.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		log.Println(err)
		return "", err
	}

	candidate := resp.Candidates[0].Content.Parts[0].(genai.Text)
	return string(candidate), nil
}
