package firestore

import (
	"context"
	"log"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

var firestoreClient *firestore.Client

func InitFirestore(ctx context.Context) {
	sa := option.WithCredentialsFile("C:\\Works\\Learning\\Golang\\golang-messenger-webhook\\internal\\firestore\\serviceAccount.json")
	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		log.Fatalln(err)
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
	}

	firestoreClient = client
}

func GetClient() *firestore.Client {
	return firestoreClient
}
