package main

import (
	"context"
	"crypto/sha256"
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

// transformEmail encrypts a string using a SHA256 hash.
// This is for deidentification of email addresses used as keys in the
// production database.
func transformEmail(email string) string {
	sha := sha256.New()
	sha.Write([]byte(email))
	encryptedEmail := fmt.Sprintf("%x", sha.Sum(nil))
	return encryptedEmail
}
