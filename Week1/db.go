/*
*
The data will be persisted in Firestore. The document structure will look like this:

Collection name: RickSteves

Subcollection: user email address

Documents
+ User query
+ Bot response
+ Timestamp

RickSteves [
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
	"log"
	"time"

	"cloud.google.com/go/firestore"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const DBName string = "l200"
const CollectionName string = "RickSteves"

type ConversationBit struct {
	BotResponse string
	UserQuery   string
	Created     time.Time
}

type ConversationHistory struct {
	UserEmail     string
	Conversations []ConversationBit
}

func saveConversation(convo ConversationBit, userEmail, projectID string) error {
	ctx := context.Background()
	client, err := firestore.NewClientWithDatabase(ctx, projectID, DBName)
	if err != nil {
		return err
	}
	defer client.Close()

	docRef := client.Collection(CollectionName).Doc(userEmail)
	docSnap, err := docRef.Get(ctx)

	if status.Code(err) == codes.NotFound {
		_, err := createHistory(docRef, userEmail, convo)
		if err != nil {
			return err
		}
		return nil
	} else if err != nil {
		return err
	}

	conversations, err := docSnap.DataAt("Conversations")
	if err != nil {
		return err
	}
	log.Printf("Retrieved conversations: \n%v\n\n", conversations)
	conversations = append(conversations, convo)

	_, err = docRef.Update(ctx, []firestore.Update{
		{
			Path:  "Conversations",
			Value: conversations,
		},
	})

	return err
}

func createHistory(docRef *firestore.DocumentRef, userEmail string, convo ConversationBit) (*firestore.DocumentSnapshot, error) {
	ctx := context.Background()

	history := &ConversationHistory{
		UserEmail:     userEmail,
		Conversations: []ConversationBit{convo},
	}

	_, err := docRef.Create(ctx, history)
	if err != nil {
		return nil, err
	}

	return docRef.Get(ctx)
}
