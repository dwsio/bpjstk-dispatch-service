package model

type ConsumedKafkaMsg struct {
	TypeId        string `json:"type_id"`
	TypeName      string `json:"type_name"`
	CategoryName  string `json:"category_name"`
	ChannelName   string `json:"channel_name"`
	PriorityOrder int    `json:"priority_order"`
	ContentHash   string `json:"hash"`
	PreferenceUrl string `json:"preference_url"`
	Data          []byte `json:"data"`
}

type PublishedKafkaMsg struct {
	TypeName      string `json:"type_name"`
	CategoryName  string `json:"category_name"`
	ChannelName   string `json:"channel_name"`
	PriorityOrder int    `json:"priority_order"`
	PreferenceUrl string `json:"preference_url"`
	Data          []byte `json:"data"`
}
