package main

import (
	"context"
	"fmt"

	dlp "cloud.google.com/go/dlp/apiv2"
	"cloud.google.com/go/dlp/apiv2/dlppb"
)

var (
	// See https://cloud.google.com/dlp/docs/infotypes-reference for complete list of info types.
	infoTypeNames = []string{"EMAIL_ADDRESS", "AGE", "FIRST_NAME", "LAST_NAME"}
)

// deidentify cleans sensitive data by replacing infoType.
//
// Taken from here: https://cloud.google.com/sensitive-data-protection/docs/samples/dlp-deidentify-replace-infotype
func deidentify(projectID, item string) (string, error) {
	ctx := context.Background()

	client, err := dlp.NewClient(ctx)
	if err != nil {
		return "", err
	}
	defer client.Close()

	input := &dlppb.ContentItem{
		DataItem: &dlppb.ContentItem_Value{
			Value: item,
		},
	}

	var infoTypes []*dlppb.InfoType
	for _, it := range infoTypeNames {
		infoTypes = append(infoTypes, &dlppb.InfoType{Name: it})
	}

	//  Associate de-identification type with info type.
	transformation := &dlppb.DeidentifyConfig_InfoTypeTransformations{
		InfoTypeTransformations: &dlppb.InfoTypeTransformations{
			Transformations: []*dlppb.InfoTypeTransformations_InfoTypeTransformation{
				{
					PrimitiveTransformation: &dlppb.PrimitiveTransformation{
						Transformation: &dlppb.PrimitiveTransformation_ReplaceWithInfoTypeConfig{},
					},
				},
			},
		},
	}

	req := &dlppb.DeidentifyContentRequest{
		Parent: fmt.Sprintf("projects/%s/locations/global", projectID),
		InspectConfig: &dlppb.InspectConfig{
			InfoTypes: infoTypes,
		},
		DeidentifyConfig: &dlppb.DeidentifyConfig{
			Transformation: transformation,
		},
		Item: input,
	}

	resp, err := client.DeidentifyContent(ctx, req)
	if err != nil {
		return "", err
	}

	return resp.GetItem().GetValue(), nil
}
