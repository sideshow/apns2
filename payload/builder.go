package payload

import "encoding/json"

type payload struct {
	content map[string]interface{}
}

type aps struct {
	Alert            interface{} `json:"alert,omitempty"`
	Badge            interface{} `json:"badge,omitempty"`
	Category         string      `json:"category,omitempty"`
	ContentAvailable int         `json:"content-available,omitempty"`
	URLArgs          []string    `json:"url-args,omitempty"`
	Sound            string      `json:"sound,omitempty"`
}

type alert struct {
	Action       string   `json:"action,omitempty"`
	ActionLocKey string   `json:"action-loc-key,omitempty"`
	Body         string   `json:"body,omitempty"`
	LaunchImage  string   `json:"launch-image,omitempty"`
	LocArgs      []string `json:"loc-args,omitempty"`
	LocKey       string   `json:"loc-key,omitempty"`
	Title        string   `json:"title,omitempty"`
	TitleLocArgs []string `json:"title-loc-args,omitempty"`
	TitleLocKey  string   `json:"title-loc-key,omitempty"`
}

func NewPayload() *payload {
	return &payload{
		map[string]interface{}{
			"aps": &aps{},
		},
	}
}

func (p *payload) Alert(alert interface{}) *payload {
	p.aps().Alert = alert
	return p
}

func (p *payload) Badge(b int) *payload {
	p.aps().Badge = b
	return p
}

func (p *payload) ZeroBadge() *payload {
	p.aps().Badge = 0
	return p
}

func (p *payload) UnsetBadge() *payload {
	p.aps().Badge = nil
	return p
}

func (p *payload) Sound(sound string) *payload {
	p.aps().Sound = sound
	return p
}

func (p *payload) ContentAvailable() *payload {
	p.aps().ContentAvailable = 1
	return p
}

func (p *payload) NewsstandAvailable() *payload {
	return p.ContentAvailable()
}

// Custom payload

func (p *payload) Custom(key string, val interface{}) *payload {
	p.content[key] = val
	return p
}

// Custom alert

func (p *payload) AlertTitle(title string) *payload {
	p.aps().alert().Title = title
	return p
}

func (p *payload) AlertTitleLocKey(key string) *payload {
	p.aps().alert().TitleLocKey = key
	return p
}

func (p *payload) AlertTitleLocArgs(args []string) *payload {
	p.aps().alert().TitleLocArgs = args
	return p
}

func (p *payload) AlertBody(body string) *payload {
	p.aps().alert().Body = body
	return p
}

func (p *payload) AlertLaunchImage(image string) *payload {
	p.aps().alert().LaunchImage = image
	return p
}

func (p *payload) AlertLocKey(key string) *payload {
	p.aps().alert().LocKey = key
	return p
}

func (p *payload) AlertAction(action string) *payload {
	p.aps().alert().Action = action
	return p
}

func (p *payload) AlertActionLocKey(key string) *payload {
	p.aps().alert().ActionLocKey = key
	return p
}

// General

func (p *payload) Category(category string) *payload {
	p.aps().Category = category
	return p
}

func (p *payload) Mdm(mdm string) *payload {
	p.content["mdm"] = mdm
	return p
}

func (p *payload) URLArgs(URLArgs []string) *payload {
	p.aps().URLArgs = URLArgs
	return p
}

func (p *payload) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.content)
}

func (p *payload) aps() *aps {
	return p.content["aps"].(*aps)
}

func (a *aps) alert() *alert {
	if _, ok := a.Alert.(*alert); !ok {
		a.Alert = &alert{}
	}
	return a.Alert.(*alert)
}
