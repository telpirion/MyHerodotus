package main

import (
	"context"
	"fmt"
	"log"
	"os/exec"
	"strings"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/cloudevents/sdk-go/v2/event"

	//"cloud.google.com/go/bigquery"
	"cloud.google.com/go/storage"
)

var (
	bucketName        = "myherodotus"
	databaseName      = "l200"
	collectionName    = "HerodotusDev"
	subCollectionname = "Conversations"
)

func init() {
	functions.CloudEvent("CollectData", collectData)
}

func main() {
	collectData(context.Background(), event.New())
}

// collectData consumes a CloudEvent message and logs details about the changed object.
func collectData(ctx context.Context, e event.Event) error {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Error: %v", r)
		}
	}()

	fmt.Println("Starting collection job")
	ok, err := exportFromFirestore(bucketName, databaseName, collectionName, subCollectionname)
	if err != nil {
		log.Println(err)
	}

	if !ok {
		log.Println("Export job failed.")
	}

	uri, err := getStorageURI(bucketName)
	fmt.Printf("Export output location: %s", uri)

	return nil
}

func getStorageURI(bucketName string) (string, error) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return "", err
	}
	defer client.Close()

	//bucket := client.Bucket(bucketName)

	return "", nil
}

// exportFromFirestore calls the gcloud command to export collections from Firestore.
// This function uses the `exec.Command()` method to invoke gcloud, which is probably
// an anti-pattern.
func exportFromFirestore(bucketName, databaseName string, collectionNames ...string) (bool, error) {

	collectionString := strings.Join(collectionNames, ",")
	cmd := exec.Command("gcloud", "firestore", "export", "gs://"+bucketName, "--database="+databaseName, "--collection-ids="+collectionString)

	stdout, err := cmd.Output()
	if err != nil {
		return false, fmt.Errorf("%s", stdout)
	}
	// Parse stdout (as YAML?)

	// Extract outputUriPrefix, name

	// Create a new longrunning operations client

	// Wait for the job to complete

	// Return the outputUriPrefix

	log.Println(string(stdout))

	return true, nil
}
