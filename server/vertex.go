package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"text/template"
	"time"

	aiplatform "cloud.google.com/go/aiplatform/apiv1"
	"cloud.google.com/go/aiplatform/apiv1/aiplatformpb"
	"cloud.google.com/go/vertexai/genai"
	"google.golang.org/api/option"
	"google.golang.org/protobuf/types/known/structpb"
)

const (
	// Minimum token count to start caching is 32768; ~110 tokens per query/response ConversationBit
	MinimumConversationNum       = 400
	GeminiTemplate               = "templates/gemini.2024.10.25.tmpl"
	GemmaTemplate                = "templates/gemma.2024.10.25.tmpl"
	GeminiModel                  = "gemini-1.5-flash-001"
	HistoryTemplate              = "templates/conversation_history.tmpl"
	MaxGemmaTokens         int32 = 2048
)

var cachedContext string = ""

type MinCacheNotReachedError struct {
	ConversationCount int
}

func (m *MinCacheNotReachedError) Error() string {
	return fmt.Sprintf("minimum context cache sized not reached; have %d, need %d",
		m.ConversationCount, MinimumConversationNum)
}

type promptInput struct {
	Query   string
	History string
}

// getTokenCount uses the Gemini tokenizer to count the tokens in some text.
func getTokenCount(text string) (int32, error) {
	location := "us-west1"
	ctx := context.Background()
	client, err := genai.NewClient(ctx, projectID, location)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	model := client.GenerativeModel(GeminiModel)

	resp, err := model.CountTokens(ctx, genai.Text(text))
	if err != nil {
		return 0, err
	}

	return resp.TotalTokens, nil
}

// setConversationContext creates string out of past conversation between user and model.
// This conversation history is used as grounding for the prompt template.
func setConversationContext(convoHistory []ConversationBit) error {
	tmp, err := template.ParseFiles(HistoryTemplate)
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	tmp.Execute(&buf, convoHistory)
	cachedContext = buf.String()
	return nil
}

// extractAnswer cleans up the response returned from the models
func extractAnswer(response string) string {
	// I am not a regex expert :/
	re := regexp.MustCompile(`### Assistant:(?s)(.*)##ENDRESPONSE##`)
	botAnswer := string(re.Find([]byte(response)))
	botAnswer = strings.Replace(botAnswer, "### Assistant:", "", 1)
	botAnswer = strings.Replace(botAnswer, "##ENDRESPONSE##", "", 1)
	if botAnswer == "" {
		botAnswer = response
	}
	return botAnswer
}

// createPrompt generates a new prompt based upon the stored prompt template.
func createPrompt(message, templateName string) (string, error) {
	tmp, err := template.ParseFiles(templateName)
	if err != nil {
		return "", nil
	}

	promptInputs := promptInput{
		Query:   message,
		History: cachedContext,
	}

	var buf bytes.Buffer
	err = tmp.Execute(&buf, promptInputs)
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

	prompt, err := createPrompt(message, GemmaTemplate)
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
func textPredictGemini(message, projectID, modelVersion string) (string, error) {
	ctx := context.Background()
	location := "us-west1"

	client, err := genai.NewClient(ctx, projectID, location)
	if err != nil {
		LogError(fmt.Sprintf("unable to create genai client: %v\n", err))
		return "", err
	}
	defer client.Close()

	modelName := GeminiModel
	if modelVersion == "gemini-tuned" {
		endpointID := os.Getenv("TUNED_MODEL_ENDPOINT_ID")
		modelName = fmt.Sprintf("projects/%s/locations/%s/endpoints/%s", projectID, location, endpointID)
	}
	llm := client.GenerativeModel(modelName)

	if convoContext != "" {
		llm.CachedContentName = convoContext
	}

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

	candidate, err := getCandidate(resp)
	if err != nil {
		LogError(err.Error())
		return "I'm not sure how to answer that. Would you please repeat the question?", nil
	}
	return extractAnswer(candidate), nil
}

// getCandidate parses the response from the model.
// It returns errors in cases where the response doesn't contain candidates
// or the candidate's parts are empty.
func getCandidate(resp *genai.GenerateContentResponse) (string, error) {
	candidates := resp.Candidates
	if len(candidates) == 0 {
		return "", errors.New("no candidates returned from model")
	}
	firstCandidate := candidates[0]
	parts := firstCandidate.Content.Parts
	if len(parts) == 0 {
		return "", errors.New("no parts in first candidate from model")
	}
	candidate := parts[0].(genai.Text)
	return string(candidate), nil
}

// storeConversationContext uploads past user conversations with the model into a Gen AI context.
// This context is used when the model is answering questions from the user.
func storeConversationContext(conversationHistory []ConversationBit, projectID string) (string, error) {
	if len(conversationHistory) < MinimumConversationNum {
		return "", &MinCacheNotReachedError{ConversationCount: len(conversationHistory)}
	}

	ctx := context.Background()
	location := "us-west1"
	client, err := genai.NewClient(ctx, projectID, location)
	if err != nil {
		return "", fmt.Errorf("unable to create client: %w", err)
	}
	defer client.Close()

	var userParts []genai.Part
	var modelParts []genai.Part
	for _, p := range conversationHistory {
		userParts = append(userParts, genai.Text(p.UserQuery))
		modelParts = append(modelParts, genai.Text(p.BotResponse))
	}

	content := &genai.CachedContent{
		Model:      GeminiModel,
		Expiration: genai.ExpireTimeOrTTL{TTL: 60 * time.Minute},
		Contents: []*genai.Content{
			{
				Role:  "user",
				Parts: userParts,
			},
			{
				Role:  "model",
				Parts: modelParts,
			},
		},
	}
	result, err := client.CreateCachedContent(ctx, content)
	if err != nil {
		return "", fmt.Errorf("CreateCachedContent: %w", err)
	}
	resourceName := result.Name

	return resourceName, nil
}
