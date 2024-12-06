/*
*
The data will be persisted in Firestore. The document structure will look like this:

Collection name: Herodotus

Subcollection: "Conversations"

Herodotus [

	{
		[EMAIL] {
			Conversations: [
				{
					BotResponse: string
					UserQuery: string
					Created: timestamp
					Model: string
					Prompt: string
					rating: string ("thumbUp" or "thumbDown")
				}
			]
	  	}
	}

]
*/
package databases

import (
	"context"
	"os"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"

	"github.com/telpirion/MyHerodotus/generated"
)

const (
	DBName            string = "l200"
	SubCollectionName string = "Conversations"
)

var CollectionName string = "HerodotusDev"

/*
	type ConversationBit struct {
		BotResponse string
		UserQuery   string
		Model       string
		Prompt      string
		Created     time.Time
		TokenCount  int32
	}
*/
type ConversationHistory struct {
	UserEmail     string
	Conversations []generated.ConversationBit
}

func SaveConversation(convo generated.ConversationBit, userEmail, projectID string) (string, error) {
	ctx := context.Background()

	// Get CollectionName for running in staging or prod
	_collectionName, ok := os.LookupEnv("COLLECTION_NAME")
	if ok {
		CollectionName = _collectionName
	}

	client, err := firestore.NewClientWithDatabase(ctx, projectID, DBName)
	if err != nil {
		return "", err
	}
	defer client.Close()

	docRef := client.Collection(CollectionName).Doc(userEmail)
	conversations := docRef.Collection(SubCollectionName)
	docRef, _, err = conversations.Add(ctx, convo)

	return docRef.ID, err
}

func UpdateConversation(documentId, userEmail, rating, projectID string) error {

	// Get CollectionName for running in staging or prod
	_collectionName, ok := os.LookupEnv("COLLECTION_NAME")
	if ok {
		CollectionName = _collectionName
	}

	ctx := context.Background()
	client, err := firestore.NewClientWithDatabase(ctx, projectID, DBName)
	if err != nil {
		return err
	}
	defer client.Close()

	docRef := client.Collection(CollectionName).Doc(userEmail).Collection(SubCollectionName).Doc(documentId)
	docRef.Set(ctx, map[string]interface{}{
		"Rating": rating,
	}, firestore.Merge(firestore.FieldPath{"Rating"}))

	return nil
}

func GetConversation(userEmail, projectID string) ([]generated.ConversationBit, error) {
	ctx := context.Background()
	conversations := []generated.ConversationBit{}
	client, err := firestore.NewClientWithDatabase(ctx, projectID, DBName)
	if err != nil {
		return conversations, err
	}
	defer client.Close()

	// Check whether this user exists in the database or not
	docRef := client.Collection(CollectionName).Doc(userEmail)
	iter := docRef.Collection(SubCollectionName).Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return conversations, err
		}
		var convo generated.ConversationBit
		err = doc.DataTo(&convo)
		if err != nil {

			continue
		}
		conversations = append(conversations, convo)
	}
	return conversations, nil
}
