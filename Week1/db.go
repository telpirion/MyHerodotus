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
	"time"

	"cloud.google.com/go/firestore"
)

const (
	DBName            string = "l200"
	CollectionName    string = "Herodotus"
	SubCollectionName string = "Conversations"
)

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
	conversations := docRef.Collection(SubCollectionName)
	_, _, err = conversations.Add(ctx, convo)

	return err
}
