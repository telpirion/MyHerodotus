package ai

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

	db "github.com/telpirion/MyHerodotus/databases"
	"github.com/telpirion/MyHerodotus/generated"
)

const (
	// Minimum token count to start caching is 32768; ~110 tokens per query/response ConversationBit
	MinimumConversationNum       = 400
	GeminiTemplate               = "templates/gemini.2024.10.25.tmpl"
	GemmaTemplate                = "templates/gemma.2024.10.25.tmpl"
	GeminiModel                  = "gemini-1.5-flash-001"
	HistoryTemplate              = "templates/conversation_history.tmpl"
	EmbeddingModelName           = "text-embedding-005"
	MaxGemmaTokens         int32 = 2048
	location                     = "us-west1"
)

var (
	cachedContext = ""
	convoContext  = ""
)

type Modality int

const (
	Gemini Modality = iota
	GeminiTuned
	Gemma
	AgentAssisted
	EmbeddingsAssisted
)

var (
	modalitiesMap = map[string]Modality{
		"gemini":              Gemini,
		"gemini-tuned":        GeminiTuned,
		"gemma":               Gemma,
		"agent-assisted":      AgentAssisted,
		"embeddings-assisted": EmbeddingsAssisted,
	}
)

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

func Predict(query, modality, projectID string) (response string, templateName string, err error) {

	switch modalitiesMap[strings.ToLower(modality)] {
	case Gemini:
		response, err = textPredictGemini(query, projectID, Gemini)
	case Gemma:
		response, err = textPredictGemma(query, projectID)
	case GeminiTuned:
		response, err = textPredictGemini(query, projectID, GeminiTuned)
	case AgentAssisted:
		response, err = textPredictWithReddit(query, projectID)
	case EmbeddingsAssisted:
		response, err = textPredictWithEmbeddings(query, projectID)
	default:
		response, err = textPredictGemini(query, projectID, Gemini)
	}

	if err != nil {
		return "", "", err
	}

	cachedContext += fmt.Sprintf("### Human: %s\n### Assistant: %s\n", query, response)
	return response, templateName, nil
}

// GetTokenCount uses the Gemini tokenizer to count the tokens in some text.
func GetTokenCount(text, projectID string) (int32, error) {
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

// SetConversationContext creates string out of past conversation between user and model.
// This conversation history is used as grounding for the prompt template.
func SetConversationContext(convoHistory []generated.ConversationBit) error {
	tmp, err := template.ParseFiles(HistoryTemplate)
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	tmp.Execute(&buf, convoHistory)
	cachedContext = buf.String()
	return nil
}

// storeConversationContext uploads past user conversations with the model into a Gen AI context.
// This context is used when the model is answering questions from the user.
func StoreConversationContext(conversationHistory []generated.ConversationBit, projectID string) (string, error) {
	if len(conversationHistory) < MinimumConversationNum {
		return "", &MinCacheNotReachedError{ConversationCount: len(conversationHistory)}
	}

	ctx := context.Background()
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
func createPrompt(message, templateName, history string) (string, error) {
	tmp, err := template.ParseFiles(templateName)
	if err != nil {
		return "", nil
	}

	promptInputs := promptInput{
		Query:   message,
		History: history,
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
	endpointID := os.Getenv("ENDPOINT_ID")
	gemma2Endpoint := fmt.Sprintf("projects/%s/locations/%s/endpoints/%s", projectID, location, endpointID)

	apiEndpoint := fmt.Sprintf("%s-aiplatform.googleapis.com:443", location)
	client, err := aiplatform.NewPredictionClient(ctx, option.WithEndpoint(apiEndpoint))
	if err != nil {
		return "", err
	}
	defer client.Close()

	parameters := map[string]interface{}{}

	prompt, err := createPrompt(message, GemmaTemplate, cachedContext)
	if err != nil {
		return "", err
	}

	tokenCount, err := GetTokenCount(prompt, projectID)
	if err != nil {
		return "", fmt.Errorf("error counting input tokens: %w", err)
	}
	if tokenCount > MaxGemmaTokens {
		prompt, err = createPrompt(message, GemmaTemplate, trimContext())
	}
	if err != nil {
		prompt = message
	}

	promptValue, err := structpb.NewValue(map[string]interface{}{
		"inputs":     prompt,
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
func textPredictGemini(message, projectID string, modality Modality) (string, error) {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, projectID, location)
	if err != nil {
		return "", err
	}
	defer client.Close()

	modelName := GeminiModel
	if modality == GeminiTuned {
		endpointID := os.Getenv("TUNED_MODEL_ENDPOINT_ID")
		modelName = fmt.Sprintf("projects/%s/locations/%s/endpoints/%s", projectID, location, endpointID)
	}
	llm := client.GenerativeModel(modelName)

	if convoContext != "" {
		llm.CachedContentName = convoContext
	}

	prompt, err := createPrompt(message, GeminiTemplate, cachedContext)
	if err != nil {
		return "", err
	}

	resp, err := llm.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return "", err
	}

	candidate, err := getCandidate(resp)
	if err != nil {
		return "", nil
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

func trimContext() (last string) {
	sep := "###"
	convos := strings.Split(cachedContext, sep)
	length := len(convos)
	if len(convos) > 3 {
		last = strings.Join(convos[length-3:length-1], sep)
	}
	return last
}

func textPredictWithReddit(query, projectID string) (string, error) {
	funcName := "GetRedditPosts"
	ctx := context.Background()
	client, err := genai.NewClient(ctx, projectID, "us-west1")
	if err != nil {
		return "", err
	}
	defer client.Close()

	schema := &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"location": {
				Type:        genai.TypeString,
				Description: "the place the user wants to go, e.g. Crete, Greece",
			},
		},
		Required: []string{"location"},
	}

	redditTool := &genai.Tool{
		FunctionDeclarations: []*genai.FunctionDeclaration{{
			Name:        funcName,
			Description: "Get Reddit posts about a location from the Travel subreddit",
			Parameters:  schema,
		}},
	}

	model := client.GenerativeModel(GeminiModel)
	model.Tools = []*genai.Tool{redditTool}

	session := model.StartChat()

	res, err := session.SendMessage(ctx, genai.Text(query))
	if err != nil {
		return "", nil
	}

	part := res.Candidates[0].Content.Parts[0]
	funcCall, ok := part.(genai.FunctionCall)
	if !ok {
		return "", fmt.Errorf("expected function call: %v", part)
	}
	if funcCall.Name != funcName {
		return "", fmt.Errorf("expected %s, got: %v", funcName, funcCall.Name)
	}
	locArg, ok := funcCall.Args["location"].(string)
	if !ok {
		return "", fmt.Errorf("expected string, got: %v", funcCall.Args["location"])
	}

	redditData, err := getRedditPosts(locArg)
	if err != nil {
		return "", err
	}

	res, err = session.SendMessage(ctx, genai.FunctionResponse{
		Name: redditTool.FunctionDeclarations[0].Name,
		Response: map[string]any{
			"output": redditData,
		},
	})
	if err != nil {
		return "", err
	}

	output := string(res.Candidates[0].Content.Parts[0].(genai.Text))
	return output, nil
}

func textPredictWithEmbeddings(query, projectID string) (string, error) {

	queryEmbed, err := getQueryTextEmbedding(query, projectID)
	if err != nil {
		return "", err
	}

	// Get context from embeddings
	embeddingContext, err := db.GetEmbedding(queryEmbed, projectID)
	if err != nil {
		return "", err
	}

	ctx := context.Background()
	client, err := genai.NewClient(ctx, projectID, location)
	if err != nil {
		return "", err
	}
	defer client.Close()

	llm := client.GenerativeModel(GeminiModel)

	createPrompt(query, GeminiTemplate, embeddingContext)

	resp, err := llm.GenerateContent(ctx, genai.Text(query))
	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}

	candidate, err := getCandidate(resp)
	if err != nil {
		return "", err
	}
	return extractAnswer(candidate), nil
}

func getQueryTextEmbedding(query, projectID string) ([]float32, error) {

	var embedding []float32
	ctx := context.Background()

	apiEndpoint := fmt.Sprintf("%s-aiplatform.googleapis.com:443", location)
	dimensionality := 128
	texts := []string{query}

	client, err := aiplatform.NewPredictionClient(ctx, option.WithEndpoint(apiEndpoint))
	if err != nil {
		return embedding, err
	}
	defer client.Close()

	endpoint := fmt.Sprintf("projects/%s/locations/%s/publishers/google/models/%s",
		projectID, location, EmbeddingModelName)
	instances := make([]*structpb.Value, len(texts))
	for i, text := range texts {
		instances[i] = structpb.NewStructValue(&structpb.Struct{
			Fields: map[string]*structpb.Value{
				"content":   structpb.NewStringValue(text),
				"task_type": structpb.NewStringValue("RETRIEVAL_QUERY"),
			},
		})
	}

	params := structpb.NewStructValue(&structpb.Struct{
		Fields: map[string]*structpb.Value{
			"outputDimensionality": structpb.NewNumberValue(float64(dimensionality)),
		},
	})

	req := &aiplatformpb.PredictRequest{
		Endpoint:   endpoint,
		Instances:  instances,
		Parameters: params,
	}
	resp, err := client.Predict(ctx, req)
	if err != nil {
		return embedding, err
	}
	embeddings := make([][]float32, len(resp.Predictions))
	for i, prediction := range resp.Predictions {
		values := prediction.GetStructValue().Fields["embeddings"].GetStructValue().Fields["values"].GetListValue().Values
		embeddings[i] = make([]float32, len(values))
		for j, value := range values {
			embeddings[i][j] = float32(value.GetNumberValue())
		}
	}

	if len(embeddings) == 0 {
		return embedding, fmt.Errorf("vertex: text embeddings: no embeddings created")
	}

	return embeddings[0], nil
}
