package firestore

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"metawebhook/internal/config"
	"metawebhook/models"
	"net/http"

	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

type ConversationsResponse struct {
	Data []DataItem `json:"data"`
}

type DataItem struct {
	ID           string `json:"id,omitempty"`
	Link         string `json:"link,omitempty"`
	UpdatedTime  string `json:"updated_time,omitempty"`
	MessageCount int    `json:"message_count,omitempty"`
	UnreadCount  int    `json:"unread_count,omitempty"`
	Messages     struct {
		Data []struct {
			From struct {
				Name  string `json:"name,omitempty"`
				Email string `json:"email,omitempty"`
				ID    string `json:"id,omitempty"`
			}
			To struct {
				Name  string `json:"name,omitempty"`
				Email string `json:"email,omitempty"`
				ID    string `json:"id,omitempty"`
			}
			CreatedTime string `json:"created_time,omitempty"`
			Message     string `json:"message,omitempty"`
			Tags        struct {
				Data []struct {
					Name string `json:"name,omitempty"`
				} `json:"data,omitempty"`
			} `json:"tags,omitempty"`
		} `json:"data,omitempty"`
	} `json:"messages,omitempty"`
}

func Save(senderID string, timestamp int, message *models.MessageType) (*models.MessageType, error) {
	// Use a service account
	ctx := context.Background()
	sa := option.WithCredentialsFile("c:\\Works\\Learning\\Golang\\golang-messenger-webhook\\internal\\firestore\\serviceAccount.json")
	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		log.Fatalln(err)
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
	}

	// Find ConversationID based on the senderID
	// Construct the API URL
	url := fmt.Sprintf("https://graph.facebook.com/v20.0/%s/conversations?user_id=%s&access_token=%s&fields=unread_count,message_count,updated_time,messages.limit(1){message,from,to,created_time,tags}", config.PageID, senderID, config.PageAccessToken)
	log.Println("Conversations API URL:", url)

	// Make the API request
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Check the status code
	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body) // Handle reading error
		return nil, fmt.Errorf("error: %s, status code: %d", body, resp.StatusCode)
	}

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	// Unmarshal the JSON response into the Response struct
	var response ConversationsResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling JSON: %w", err)
	}

	conversationID := response.Data[0].ID
	messageCount := response.Data[0].MessageCount
	unreadCount := response.Data[0].UnreadCount
	updatedTime := response.Data[0].UpdatedTime
	messageCreatedTime := response.Data[0].Messages.Data[0].CreatedTime
	messageFrom := response.Data[0].Messages.Data[0].From

	messageID := message.Mid
	messageText := message.Text
	tags := response.Data[0].Messages.Data[0].Tags.Data

	inbox := false
	platform := "messenger"
	read := false
	sent := false
	source := "web"
	for _, tag := range tags {
		switch tag.Name {
		case "inbox":
			inbox = true
		case "read":
			read = true
		case "sent":
			sent = true
		case "messenger":
			platform = "messenger"
		case "source:mobile":
			source = "mobile"
		}
	}

	// Save Message information in Conversations Collection
	_, err = firestoreClient.Collection("conversations").Doc(conversationID).Set(ctx, map[string]interface{}{
		"id":            conversationID,
		"message_count": messageCount,
		"unread_count":  unreadCount,
		"updated_time":  updatedTime,
		"last_message": map[string]interface{}{
			"data": []interface{}{
				map[string]interface{}{
					"id":           messageID,
					"created_time": messageCreatedTime,
					"from": map[string]interface{}{
						"email": messageFrom.Email,
						"id":    messageFrom.ID,
						"name":  messageFrom.Name,
					},
					"message":  messageText,
					"read":     read,
					"sent":     sent,
					"platform": platform,
					"source":   source,
					"inbox":    inbox,
				},
			},
		},
	})
	if err != nil {
		log.Fatalf("Failed adding a new message: %v", err)
		return nil, err
	}

	// Save Message information in Messages Collection
	_, err = client.Collection("conversations").Doc(conversationID).Collection("messages").Doc(messageID).Set(ctx, map[string]interface{}{
		"id":           messageID,
		"created_time": messageCreatedTime,
		"from": map[string]interface{}{
			"email": messageFrom.Email,
			"id":    messageFrom.ID,
			"name":  messageFrom.Name,
		},
		"message":  messageText,
		"read":     read,
		"sent":     sent,
		"platform": platform,
		"source":   source,
		"inbox":    inbox,
	})
	if err != nil {
		log.Fatalf("Failed adding a new message: %v", err)
		return nil, err
	}
	return message, nil
}
