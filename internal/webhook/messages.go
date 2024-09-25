package webhook

import (
	"metawebhook/internal/firestore"
	"metawebhook/models"
)

func handleMessage(senderID string, timestamp int, message *models.MessageType) {
	firestore.Save(senderID, timestamp, message)
}
