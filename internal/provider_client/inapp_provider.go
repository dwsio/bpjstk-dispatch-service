package provider_client

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	"git.bpjsketenagakerjaan.go.id/centralized-notification-system/dispatch-service/pkg/constants"
	"github.com/pkg/errors"
)

type SendInAppReq struct {
	InAppMessage InAppMessage `json:"in_app_message"`
	IsDraft      bool         `json:"is_draft"`
}

type InAppMessage struct {
	Name              string   `json:"name"`
	Location          string   `json:"location"`
	RemoveMargin      bool     `json:"remove_margin"`
	Contents          Contents `json:"contents"`
	RedisplayDelay    string   `json:"redisplay_delay"`
	RedisplayLimit    string   `json:"redisplay_limit"`
	IsTemplate        bool     `json:"is_template"`
	DisplayDuration   int      `json:"display_duration"`
	IncludeSegmentIds []string `json:"include_segment_ids"`
	ExcludeSegmentIds []string `json:"exclude_segment_ids"`
	Triggers          []string `json:"triggers"`
	StartTime         string   `json:"start_time"`
	EndTime           string   `json:"end_time"`
	IsCustomHtml      bool     `json:"is_custom_html"`
	Html              string   `json:"html"`
}

type Contents struct {
	Version string `json:"version"`
	Pages   []Page `json:"pages"`
}

type Page struct {
	Background Background `json:"background"`
	Id         string     `json:"id"`
	Elements   []Element  `json:"elements"`
}

type Background struct {
	Color string `json:"color"`
	Image string `json:"image"`
}

type Element struct {
	Font            Font   `json:"font"`
	Id              string `json:"id"`
	Text            string `json:"text"`
	Type            string `json:"type"`
	Margin          Margin `json:"margin"`
	Action          Action `json:"action"`
	Url             string `json:"url"`
	BackgroundColor string `json:"background_color"`
	BorderColor     string `json:"border_color"`
	BorderRadius    int    `json:"border_radius"`
	BorderWidth     int    `json:"border_width"`
	BoxShadowBlur   int    `json:"box_shadow_blur"`
	BoxShadowColor  string `json:"box_shadow_color"`
	BoxShadowSpread int    `json:"box_shadow_spread"`
	BoxShadowX      int    `json:"box_shadow_x"`
	BoxShadowY      int    `json:"box_shadow_y"`
	IsVisible       bool   `json:"is_visible"`
	Height          int    `json:"height"`
	Width           int    `json:"width"`
}

type Font struct {
	Size       string `json:"size"`
	Color      string `json:"color"`
	Alignment  string `json:"alignment"`
	Family     string `json:"family"`
	Decoration string `json:"decoration"`
	Style      string `json:"style"`
	Weight     string `json:"weight"`
}

type Margin struct {
	Bottom int `json:"bottom"`
	Left   int `json:"left"`
	Right  int `json:"right"`
	Top    int `json:"top"`
}

type Action struct {
	Close     bool     `json:"close"`
	Prompts   []string `json:"prompts,omitempty"`
	UrlTarget string   `json:"url_target"`
}

type SendInAppRes struct {
	Success bool    `json:"success"`
	Payload Payload `json:"payload"`
}

type Payload struct {
	Success bool `json:"success"`
}

func (pc *ProviderClient) SendInApp(ctx context.Context, url string) (string, error) {
	_, span := pc.tracer.Start(ctx, "ProviderClient.SendInApp")
	defer span.End()

	client := pc.NewHttpClient()

	req := getInAppRequest()
	httpReqBody, err := json.Marshal(req)
	if err != nil {
		return "", errors.Wrap(err, "json.Marshal")
	}

	httpReq, err := http.NewRequest(constants.METHOD_POST, url, bytes.NewBuffer(httpReqBody))
	if err != nil {
		return "", errors.Wrap(err, "http.NewRequest")
	}

	httpReq.Header.Add("Content-Type", "application/json")

	httpRes, err := client.Do(httpReq)
	if err != nil {
		return "", errors.Wrap(err, "client.Do")
	}
	defer httpRes.Body.Close()

	if httpRes.StatusCode != http.StatusOK {
		return "", errors.New("request failed")
	}

	var httpResData SendInAppRes
	err = json.NewDecoder(httpRes.Body).Decode(&httpResData)
	if err != nil {
		return "", errors.Wrap(err, "jsonDecoder.Decode")
	}

	result := "failed"
	if httpResData.Payload.Success {
		result = "success"
	}

	return result, nil
}

func getInAppRequest() *SendInAppReq {
	return &SendInAppReq{
		InAppMessage: InAppMessage{
			Name:         "Test",
			Location:     "center_modal",
			RemoveMargin: false,
			Contents: Contents{
				Version: "3",
				Pages: []Page{
					{
						Background: Background{
							Color: "#FFFFFF",
							Image: "",
						},
						Id: "efe129f9-243c-4af8-82f9-d09ad2d7c6ba",
						Elements: []Element{
							{
								Font: Font{
									Size:       "24",
									Color:      "#22",
									Alignment:  "center",
									Family:     "System Font (Default)",
									Decoration: "",
									Style:      "",
									Weight:     "",
								},
								Id:   "e5198a3f-ca91-4d1b-a9ef-51466880279a",
								Text: "Heading",
								Type: "title",
								Margin: Margin{
									Top:    0,
									Right:  0,
									Bottom: 0,
									Left:   0,
								},
							},
							{
								Id: "11703a8a-a16f-42f9-bf9f-083ecd7c2d58",
								Action: Action{
									Close:     false,
									Prompts:   []string{},
									UrlTarget: "browser",
								},
								Type: "image",
								Margin: Margin{
									Top:    0,
									Right:  0,
									Bottom: 0,
									Left:   0,
								},
								Url: "https://media.onesignal.com/iam/default_image_20200320.png",
							},
							{
								Action: Action{
									Close:     false,
									Prompts:   []string{},
									UrlTarget: "browser",
								},
								BackgroundColor: "#1F8FEB",
								BorderColor:     "#051B2C",
								BorderRadius:    4,
								BorderWidth:     0,
								BoxShadowBlur:   0,
								BoxShadowColor:  "#051B2C",
								BoxShadowSpread: 0,
								BoxShadowX:      0,
								BoxShadowY:      0,
								Font: Font{
									Size:       "24",
									Color:      "#fff",
									Alignment:  "center",
									Family:     "System Font (Default)",
									Decoration: "",
									Style:      "",
									Weight:     "",
								},
								Id:   "b46b8819-4b92-4e21-bda4-50172358f178",
								Text: "Click Me",
								Type: "button",
								Margin: Margin{
									Top:    0,
									Right:  0,
									Bottom: 0,
									Left:   0,
								},
							},
							{
								Id:        "ecb00f21-c2a5-4925-9449-8251a85639b2",
								Type:      "close_button",
								IsVisible: true,
								Height:    10,
								Width:     10,
								Action: Action{
									Close:     false,
									Prompts:   []string{},
									UrlTarget: "browser",
								},
								Margin: Margin{
									Top:    0,
									Right:  0,
									Bottom: 0,
									Left:   0,
								},
							},
							{
								Type: "body",
								Action: Action{
									Close:     false,
									Prompts:   []string{},
									UrlTarget: "browser",
								},
								Id: "df7188e9-d92d-4a46-a596-40674200b11b",
							},
						},
					},
				},
			},
			RedisplayLimit:    "",
			RedisplayDelay:    "",
			IsTemplate:        false,
			DisplayDuration:   0,
			IncludeSegmentIds: []string{},
			ExcludeSegmentIds: []string{},
			Triggers:          []string{},
			StartTime:         "2023-07-17T03:38:15.197Z",
			EndTime:           "",
			IsCustomHtml:      false,
			Html:              "",
		},
		IsDraft: false,
	}
}
