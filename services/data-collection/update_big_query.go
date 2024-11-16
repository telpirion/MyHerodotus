package datacollection

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/cloudevents/sdk-go/v2/event"
	"google.golang.org/protobuf/encoding/protojson"

	"cloud.google.com/go/bigquery"
	"cloud.google.com/go/firestore"
	"github.com/googleapis/google-cloudevents-go/cloud/firestoredata"
)

var (
	projectID      = ""
	datasetName    = ""             // Pass in as env var
	collectionName = "HerodotusDev" // Pass in as env var
)

const (
	databaseID        = "l200"
	tableName         = "conversations"
	subCollectionname = "Conversations"
)

func init() {
	functions.CloudEvent("CollectData", collectData)
}

// collectData consumes a CloudEvent message and logs details about the changed object.
func collectData(ctx context.Context, e event.Event) error {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Error: %v", r)
		}
	}()

	projectID = os.Getenv("PROJECT_ID")
	datasetName = os.Getenv("DATASET_NAME")
	collectionName = os.Getenv("BUILD_VER")

	var data firestoredata.DocumentEventData
	options := protojson.UnmarshalOptions{
		DiscardUnknown: true,
	}
	if err := options.Unmarshal(e.Data(), &data); err != nil {
		log.Fatalf("protojson.Unmarshal: %v", err)
		return err
	}

	documentName := data.GetValue().Name

	fmt.Println("Starting collection job")
	docSnap, err := retrieveDocument(documentName)
	if err != nil {
		log.Fatalf("retrieve Firestore doc: %v", err)
		return err
	}

	ok, err := updateBigQuery(docSnap)
	if err != nil {
		log.Fatalf("update BigQuery: %v", err)
		return err
	}

	if !ok {
		log.Printf("update transaction failed: see if Rating is populated: %s", documentName)
		return nil
	}

	return nil
}

func updateBigQuery(docSnap *firestore.DocumentSnapshot) (bool, error) {
	var convo ConversationBit
	err := docSnap.DataTo(&convo)
	if err != nil {
		return false, err
	}

	// Don't update if there is no user rating.
	if convo.GetRating() == "" {
		return false, nil
	}

	ctx := context.Background()
	client, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		return false, err
	}
	defer client.Close()

	// Update the appropriate row in BigQuery
	ins := client.Dataset(datasetName).Table(tableName).Inserter()
	ins.IgnoreUnknownValues = true
	ins.SkipInvalidRows = true

	err = ins.Put(ctx, &convo)
	if err != nil {
		return false, err
	}

	return true, nil
}

// retrieveDocument gets a DocumentSnapshot from Firestore by the Document's fully-qualified name.
func retrieveDocument(documentName string) (doc *firestore.DocumentSnapshot, err error) {
	ctx := context.Background()
	client, err := firestore.NewClientWithDatabase(ctx, projectID, databaseID)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	// The format of the document name returned from EventArc is
	// projects/{proj}/databses/{db}/documents/{document_path}
	paths := strings.Split(documentName, "documents/")
	if len(paths) != 2 {
		return nil, fmt.Errorf("bad document resource name received: %s", documentName)
	}

	documentPath := paths[len(paths)-1]

	// Get the document from Firestore
	docRef := client.Doc(documentPath)

	doc, err = docRef.Get(ctx)
	if !doc.Exists() {
		return nil, fmt.Errorf("document doesn't exist")
	}

	return doc, nil
}
