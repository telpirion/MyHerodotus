package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	aiplatform "cloud.google.com/go/aiplatform/apiv1"
	"cloud.google.com/go/aiplatform/apiv1/aiplatformpb"
	"cloud.google.com/go/vertexai/genai"
	"google.golang.org/api/option"
	"google.golang.org/protobuf/types/known/structpb"
)

func extractAnswer(response string) string {
	// I am not a regex expert :/
	re := regexp.MustCompile(`##RESPONSE##(?s)(.*)##ENDRESPONSE##`)
	botAnswer := string(re.Find([]byte(response)))
	botAnswer = strings.Replace(botAnswer, "##RESPONSE##", "", 1)
	botAnswer = strings.Replace(botAnswer, "##ENDRESPONSE##", "", 1)
	return botAnswer
}

func createPrompt(message string) string {
	return `
You are a helpful travel agent. The user wants to have a conversation with you about where to go.
You are going to help them plan their trip.

Be sure to label your response with ##RESPONSE## and end your response with ##ENDRESPONSE##.

Do not include system instructions in the response.

Example:
[USER]: I want to go to Italy.
##RESPONSE## Italy is a fantastic place to go! What would you like to experience: the food, the
history, the art, the people, or some combination?##ENDRESPONSE##

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

	parameters := map[string]interface{}{}

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
	log.Println(prediction)
	value := prediction[0].GetStringValue()

	return extractAnswer(value), nil
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
