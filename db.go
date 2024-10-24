/*
*
The data will be persisted in Firestore. The document structure will look like this:

Collection name: Herodotus

Subcollection: user email address

Documents
+ User query
+ Bot response
+ Timestamp

Herodotus [

	{
		email {
			Conversations: [
				{
					BotResponse: string
					UserQuery: string
					Created: timestamp
				}
			]
	  	}
	}

]
*/
package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	DBName            string = "l200"
	SubCollectionName string = "Conversations"
)

var CollectionName string = "HerodotusDev"

type ConversationBit struct {
	BotResponse string
	UserQuery   string
	Model       string
	Prompt      string
	Created     time.Time
}

type ConversationHistory struct {
	UserEmail     string
	Conversations []ConversationBit
}

func saveConversation(convo ConversationBit, userEmail, projectID string) (string, error) {
	ctx := context.Background()

	// Get CollectionName for running in staging or prod
	_collectionName, ok := os.LookupEnv("COLLECTION_NAME")
	if ok {
		CollectionName = _collectionName
	}

	client, err := firestore.NewClientWithDatabase(ctx, projectID, DBName)
	if err != nil {
		LogError(fmt.Sprintf("firestore.Client: %v\n", err))
		return "", err
	}
	defer client.Close()

	docRef := client.Collection(CollectionName).Doc(userEmail)
	conversations := docRef.Collection(SubCollectionName)
	docRef, _, err = conversations.Add(ctx, convo)

	return docRef.ID, err
}

func updateConversation(documentId, userEmail, rating, projectID string) error {

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
		"rating": rating,
	}, firestore.Merge(firestore.FieldPath{"rating"}))

	return nil
}

func getConversation(userEmail, projectID string) ([]ConversationBit, error) {
	ctx := context.Background()
	conversations := []ConversationBit{}
	client, err := firestore.NewClientWithDatabase(ctx, projectID, DBName)
	if err != nil {
		LogError(fmt.Sprintf("firestore.Client: %v\n", err))
		return conversations, err
	}
	defer client.Close()

	// Check whether this user exists in the database or not
	docRef := client.Collection(CollectionName).Doc(userEmail)
	_, err = docRef.Get(ctx)
	if status.Code(err) == codes.NotFound {
		return conversations, nil
	}
	if err != nil {
		LogError(fmt.Sprintf("firestore.DocumentRef: %v\n", err))
		return conversations, err
	}

	iter := docRef.Collection(SubCollectionName).Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			LogError(fmt.Sprintf("Firestore Iterator: %v\n", err))
			return conversations, err
		}
		var convo ConversationBit
		err = doc.DataTo(&convo)
		if err != nil {
			LogError(fmt.Sprintf("Firestore document unmarshaling: %v\n", err))
			continue
		}
		conversations = append(conversations, convo)
	}
	return conversations, nil
}
