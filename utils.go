package main

import (
	"context"
	"fmt"

	"cloud.google.com/go/storage"
)

// Function to print a list of blobs in a Google Cloud Storage bucket
func listBlobs() {
	// ...
	// Instantiate a Cloud Storage client

	ctx := context.Background()
	storageClient, err := storage.NewClient(ctx)
	if err != nil {
		fmt.Printf("Failed to create client: %v", err)
		return
	}

	bucketName := "erschmid-test-291318-bucket"

	blobs := storageClient.Bucket(bucketName).Objects(ctx, nil)
	for {
		blob, err := blobs.Next()
		if err != nil {
			break
		}
		fmt.Println(blob.Name)
	}
	// List blobs in the bucket
	// ...
	// Print the list of blobs
	// ...
}
