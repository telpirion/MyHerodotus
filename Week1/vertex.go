package main

// [START aiplatform_text_predictions]
// [START generativeaionvertexai_text_predictions]

import (
	"context"
	"fmt"
	"log"

	aiplatform "cloud.google.com/go/aiplatform/apiv1beta1"
	"cloud.google.com/go/aiplatform/apiv1beta1/aiplatformpb"
	"google.golang.org/api/option"
	"google.golang.org/protobuf/types/known/structpb"
)

// textPredict generates text with certain prompt and configurations.
func textPredict(message, projectID, location, model string) (string, error) {
	ctx := context.Background()

	prompt := `The user wants to go on a trip somewhere. They have asked you this question:` + message
	log.Println(prompt)

	publisher := "google"
	parameters := map[string]interface{}{
		"temperature":     0.8,
		"maxOutputTokens": 256,
		"topP":            0.4,
		"topK":            40,
	}

	apiEndpoint := fmt.Sprintf("%s-aiplatform.googleapis.com:443", location)

	client, err := aiplatform.NewPredictionClient(ctx, option.WithEndpoint(apiEndpoint))
	if err != nil {
		log.Println(err)
		return "", err
	}
	defer client.Close()

	base := fmt.Sprintf("projects/%s/locations/%s/publishers/%s/models", projectID, location, publisher)
	url := fmt.Sprintf("%s/%s", base, model)

	// Instances: the prompt to use with the text model
	promptValue, err := structpb.NewValue(map[string]interface{}{
		"prompt": prompt,
	})
	if err != nil {
		log.Println(err)
		return "", err
	}

	// Parameters: the model configuration parameters
	parametersValue, err := structpb.NewValue(parameters)
	if err != nil {
		log.Println(err)
		return "", err
	}

	// PredictRequest: create the model prediction request
	req := &aiplatformpb.PredictRequest{
		Endpoint:   url,
		Instances:  []*structpb.Value{promptValue},
		Parameters: parametersValue,
	}

	// PredictResponse: receive the response from the model
	resp, err := client.Predict(ctx, req)
	if err != nil {
		log.Println(err)
		return "", err
	}

	response := resp.GetPredictions()[0].GetStringValue()

	log.Println(response)

	return response, nil
}

// [END aiplatform_text_predictions]
// [END generativeaionvertexai_text_predictions]
