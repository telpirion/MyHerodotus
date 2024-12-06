package databases

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/telpirion/MyHerodotus/generated"
)

var (
	email           string = "unittest@example.com"
	email2          string = "unittest2@example.com"
	_projectID      string
	_collectionName string = "TestCollection"
)

func TestMain(m *testing.M) {
	_projectID = os.Getenv("PROJECT_ID")
	ctx := context.Background()
	client, err := firestore.NewClientWithDatabase(ctx, _projectID, DBName)
	if err != nil {
		log.Fatal(err)

	}
	defer client.Close()

	collection := client.Collection(_collectionName)
	subcollection := collection.Doc(email2).Collection(SubCollectionName)
	_, _, err = subcollection.Add(ctx, generated.ConversationBit{
		UserQuery:   "test user query",
		BotResponse: "test bot response",
		Created:     time.Now().Unix(),
	})
	if err != nil {
		log.Fatal(err)
	}
	_, _, err = subcollection.Add(ctx, generated.ConversationBit{
		UserQuery:   "test user query 2",
		BotResponse: "test bot response 2",
		Created:     time.Now().Unix(),
	})
	if err != nil {
		log.Fatal(err)
	}

	m.Run()

	//Clean up after tests run
	_, err = collection.Doc(email).Delete(ctx)
	if err != nil {
		log.Fatal(err)
	}

	_, err = collection.Doc(email2).Delete(ctx)
	if err != nil {
		log.Fatal(err)
	}
}

func TestSaveConversation(t *testing.T) {
	convo := &generated.ConversationBit{
		UserQuery:   "This is from unit test",
		BotResponse: "This is a bot response",
		Created:     time.Now().Unix(),
	}
	id, err := SaveConversation(*convo, email, _projectID)
	if err != nil {
		t.Fatal(err)
	}

	if id == "" {
		t.Error("Empty document ID")
	}

	nextConvo := &generated.ConversationBit{
		UserQuery:   "This is also from a unit test",
		BotResponse: "This is another fake bot response",
		Created:     time.Now().Unix(),
	}
	nextID, err := SaveConversation(*nextConvo, email, _projectID)
	if err != nil {
		t.Fatalf("Error on adding next conversation: %v\n\n", err)
	}
	if nextID == "" {
		t.Error("Empty second document ID")
	}
}

func TestGetConversation(t *testing.T) {
	conversations, err := GetConversation(email2, _projectID)
	if err != nil {
		t.Fatal(err)
	}

	want := 2
	if got := len(conversations); got != want {
		t.Errorf("wanted %d messages, got %d", want, got)
	}
}
