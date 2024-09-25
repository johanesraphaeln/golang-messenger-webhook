package models

type SenderType struct {
	ID string `json:"id"`
}

type MessageType struct {
	Mid         string `json:"mid,omitempty"`
	Text        string `json:"text,omitempty"`
	Attachments []struct {
		Payload struct {
			URL string `json:"url"`
		} `json:"payload"`
	} `json:"attachments,omitempty"`
}

type PostbackType struct {
	Payload string `json:"payload"`
}
