package model

type Email struct {
	RecipientTo  []string `json:"recipient_to"`
	RecipientCc  []string `json:"recipient_cc,omitempty"`
	RecipientBcc []string `json:"recipient_bcc,omitempty"`
	Subject      string   `json:"subject,omitempty"`
	ContentText  string   `json:"content,omitempty"`
	IsHTML       bool     `json:"is_html"`
	ContentHTML  string   `json:"content_html,omitempty"`
	IsAttach     bool     `json:"is_attach"`
	Attachment   []string `json:"attachment,omitempty"`
	AttachName   []string `json:"attach_name,omitempty"`
	Status       string   `json:"status,omitempty"`
}
