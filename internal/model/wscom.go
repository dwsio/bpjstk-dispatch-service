package model

type SendSmsResCode struct {
	Account string `json:"acc"`
	Type    string `json:"type"`
	Length  int    `json:"length"`
	NumSms  int    `json:"numSMS"`
	Code    string `json:"code"`
	Status  string `json:"status"`
	Message string `json:"message"`
	MsgID   string `json:"msgid"`
}
