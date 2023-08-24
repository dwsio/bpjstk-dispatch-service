package model

type Sms struct {
	RecipientPhoneNumber string `json:"recipient_phone_number"`
	Content              string `json:"content"`
	MessageId            string `json:"message_id,omitempty"`
	Status               string `json:"status,omitempty"`
}
