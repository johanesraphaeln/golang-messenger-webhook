package main

import (
	"context"
	"log"
	"metawebhook/internal/config"
	"metawebhook/internal/firestore"
	"metawebhook/internal/webhook"
	"net/http"
)

func main() {
	ctx := context.Background()
	config.LoadConfig()
	firestore.InitFirestore(ctx)

	http.HandleFunc("/webhook", webhook.WebhookHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
