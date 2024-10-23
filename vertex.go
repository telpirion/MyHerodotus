package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"text/template"

	aiplatform "cloud.google.com/go/aiplatform/apiv1"
	"cloud.google.com/go/aiplatform/apiv1/aiplatformpb"
	"cloud.google.com/go/vertexai/genai"
	"google.golang.org/api/option"
	"google.golang.org/protobuf/types/known/structpb"
)

const (
	GeminiTemplate = "templates/gemini.tmpl"
	GemmaTemplate  = "templates/gemma.impl"
)

func extractAnswer(response string) string {
	// I am not a regex expert :/
	re := regexp.MustCompile(`##RESPONSE##(?s)(.*)##ENDRESPONSE##`)
	botAnswer := string(re.Find([]byte(response)))
	botAnswer = strings.Replace(botAnswer, "##RESPONSE##", "", 1)
	botAnswer = strings.Replace(botAnswer, "##ENDRESPONSE##", "", 1)
	if botAnswer == "" {
		botAnswer = response
	}
	return botAnswer
}

func createPrompt(message, templateName string) (string, error) {
	tmp, err := template.ParseFiles(templateName)
	if err != nil {
		return "", nil
	}
	var buf bytes.Buffer
	err = tmp.ExecuteTemplate(&buf, "gemini.tmpl", struct{ Query string }{Query: message})
	if err != nil {
		return "", err
	}
	return buf.String(), nil
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
		LogError(fmt.Sprintf("unable to create prediction client: %v\n", err))
		return "", err
	}
	defer client.Close()

	parameters := map[string]interface{}{}

	prompt, err := createPrompt(message, GeminiTemplate)
	if err != nil {
		LogError(fmt.Sprintf("unable to create Gemma prompt: %v\n", err))
		return "", err
	}

	promptValue, err := structpb.NewValue(map[string]interface{}{
		"inputs":     prompt,
		"parameters": parameters,
	})
	if err != nil {
		LogError(fmt.Sprintf("unable to create prompt value: %v\n", err))
		return "", err
	}

	req := &aiplatformpb.PredictRequest{
		Endpoint:  gemma2Endpoint,
		Instances: []*structpb.Value{promptValue},
	}

	resp, err := client.Predict(ctx, req)
	if err != nil {
		LogError(fmt.Sprintf("unable to make prediction: %v\n", err))
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
		LogError(fmt.Sprintf("unable to create genai client: %v\n", err))
		return "", err
	}
	defer client.Close()

	llm := client.GenerativeModel(model)
	prompt, err := createPrompt(message, GeminiTemplate)
	if err != nil {
		LogError(fmt.Sprintf("unable to create Gemini prompt: %v\n", err))
		return "", err
	}

	resp, err := llm.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		LogError(fmt.Sprintf("unable to generate content: %v\n", err))
		return "", err
	}

	candidate := resp.Candidates[0].Content.Parts[0].(genai.Text)
	return extractAnswer(string(candidate)), nil
}
