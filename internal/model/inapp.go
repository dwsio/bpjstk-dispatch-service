package model

type InApp struct {
	PlayerIds  []string `json:"player_ids"`
	Segments   []string `json:"segments"`
	Heading    string   `json:"heading"`
	Content    string   `json:"content"`
	PictureUrl string   `json:"picture_url"`
	IsIos      bool     `json:"is_ios"`
}
