package webhook

type TextMessage struct {
	Mid         string `json:"mid"`
	Text        string `json:"text"`
	Quick_Reply struct {
		Payload string `json:"payload"`
	} `json:"quick_reply"`
}

type ReplyMessage struct {
	Mid      string `json:"mid"`
	Text     string `json:"text"`
	Reply_To struct {
		Mid string `json:"mid"`
	} `json:"quick_reply"`
}

type AttachmentMessage struct {
	Mid         string `json:"mid"`
	Attachments struct {
		Type    string `json:"type"`
		Payload struct {
			Url string `json:"url"`
		} `json:"payload"`
	} `json:"attachments"`
}
