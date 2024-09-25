package webhook

import (
	"encoding/json"
	"log"
	"metawebhook/models"
	"net/http"
)

func WebhookHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		handleVerification(w, r)
	case http.MethodPost:
		log.Print("POST")
		handleMessages(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func handleVerification(w http.ResponseWriter, r *http.Request) {
	mode := r.URL.Query().Get("hub.mode")
	token := r.URL.Query().Get("hub.verify_token")
	challenge := r.URL.Query().Get("hub.challenge")

	if mode == "subscribe" && token == "abclimadasar" {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(challenge))
		log.Println("WEBHOOK_VERIFIED")
	} else {
		w.WriteHeader(http.StatusForbidden)
	}
}

func handleMessages(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Object string `json:"object"`
		Entry  []struct {
			ID        string `json:"id"`
			Time      int    `json:"time"`
			Messaging []struct {
				Sender    *models.SenderType   `json:"sender,omitempty"`
				Timestamp int                  `json:"timestamp,omitempty"`
				Message   *models.MessageType  `json:"message,omitempty"`
				Postback  *models.PostbackType `json:"postback,omitempty"`
			} `json:"messaging"`
		} `json:"entry"`
	}

	// Call the function that processes messages
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		log.Println("Error decoding JSON:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if body.Object == "page" {
		for _, entry := range body.Entry {
			for _, messagingEvent := range entry.Messaging {
				senderID := messagingEvent.Sender.ID
				timestamp := messagingEvent.Timestamp
				if messagingEvent.Message != nil {
					handleMessage(senderID, timestamp, (*models.MessageType)(messagingEvent.Message))
				}
			}
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("EVENT_RECEIVED"))
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}
