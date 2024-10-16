package main

import (
	"testing"
	"time"
)

var (
	email      string = "unittest@example.com"
	_projectID string = "erschmid-test-291318"
)

func TestMain(m *testing.M) {
	m.Run()

	// Clean up after tests run
	// ctx := context.Background()
	// client, err := firestore.NewClientWithDatabase(ctx, _projectID, DBName)
	// if err != nil {
	// 	log.Fatal(err)

	// }
	// collection := client.Collection(CollectionName)
	// _, err = collection.Doc(email).Delete(ctx)
	// if err != nil {
	// 	log.Fatal(err)
	// }
}

func TestSaveConversation(t *testing.T) {
	convo := &ConversationBit{
		UserQuery:   "This is from unit test",
		BotResponse: "This is a bot response",
		Created:     time.Now(),
	}
	err := saveConversation(*convo, email, _projectID)
	if err != nil {
		t.Fatal(err)
	}

	nextConvo := &ConversationBit{
		UserQuery:   "This is also from a unit test",
		BotResponse: "This is another fake bot response",
		Created:     time.Now(),
	}
	err = saveConversation(*nextConvo, email, _projectID)
	if err != nil {
		t.Fatalf("Error on adding next conversation: %v\n\n", err)
	}
}
