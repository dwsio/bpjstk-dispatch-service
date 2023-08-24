package model

type Push struct {
	PlayerIds  []string               `json:"player_ids"`
	Heading    string                 `json:"heading"`
	Content    string                 `json:"content"`
	PictureUrl string                 `json:"picture_url"`
	Data       map[string]interface{} `json:"data"`
	IsIos      bool                   `json:"is_ios"`
	MessageId  string                 `json:"message_id,omitempty"`
	Status     string                 `json:"status,omitempty"`
}
