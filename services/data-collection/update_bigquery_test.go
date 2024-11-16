package datacollection

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/googleapis/google-cloudevents-go/cloud/firestoredata"
	"google.golang.org/protobuf/encoding/protojson"
)

var (
	testProjectID  = ""
	testDatabaseID = "l200"
)

const (
	testResource = "/HerodotusDev/34ec0862ba25c6705faa0084193cb2f5ded5709569998cda5ba518b15710647a/Conversations/6QJ1lQNpf9YUBERWtwR5"
)

func TestMain(m *testing.M) {
	testProjectID = os.Getenv("PROJECT_ID")
	m.Run()
}

func TestCollectData(t *testing.T) {
	resourceName := fmt.Sprintf("projects/%s/databases/%s/documents", testProjectID, testDatabaseID)
	documentName := fmt.Sprintf("%s%s", resourceName, testResource)
	ctx := context.Background()
	evt := event.New()
	data := &firestoredata.DocumentEventData{
		Value: &firestoredata.Document{
			Name: documentName,
		},
	}
	options := protojson.MarshalOptions{
		AllowPartial: true,
	}
	encodedData, err := options.Marshal(data)
	if err != nil {
		t.Fatal(err)
	}

	evt.SetData(*event.StringOfApplicationJSON(), encodedData)

	err = collectData(ctx, evt)
	if err != nil {
		t.Fatal(err)
	}
}
