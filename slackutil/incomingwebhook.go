package slackutil

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/tsaikd/KDGoLib/errutil"
	"github.com/tsaikd/KDGoLib/webutil"
)

// errors
var (
	ErrorSlackUnexpectedResponseBody1   = errutil.NewFactory("unexpected slack response body: %v")
	ErrorSlackUnexpectedResponseStatus1 = errutil.NewFactory("unexpected slack response status: %v")
)

type IncomingWebHook struct {
	WebHookURL string `json:"webhookurl,omitempty"`
	Channel    string `json:"channel,omitempty"`
	IconEmoji  string `json:"icon_emoji,omitempty"`
}

func (t *IncomingWebHook) SendMessage(message IncomingWebHookMessage) (err error) {
	data, err := json.Marshal(message)
	if err != nil {
		return
	}
	res, err := http.Post(t.WebHookURL, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return
	}

	if err = webutil.NewResponseError(res); err != nil {
		return ErrorSlackUnexpectedResponseStatus1.New(err, res.StatusCode)
	}

	body, err := webutil.ReadResponse(res)
	if err != nil {
		return
	}

	if body != "ok" {
		return ErrorSlackUnexpectedResponseBody1.New(nil, body)
	}

	return
}

type IncomingWebHookMessage struct {
	Channel   string `json:"channel,omitempty"`
	Username  string `json:"username,omitempty"`
	Text      string `json:"text,omitempty"`
	IconEmoji string `json:"icon_emoji,omitempty"`
}

func (t *IncomingWebHookMessage) FillEmptyFieldWithDefaultValue(config IncomingWebHook) {
	if t.Channel == "" {
		t.Channel = config.Channel
		if t.Channel == "" {
			t.Channel = "#general"
		}
	}
	if t.Username == "" {
		t.Username = "milkr slack bot"
	}
	if t.Text == "" {
		t.Text = "milkr slack bot test trigger message"
	}
	if t.IconEmoji == "" {
		t.IconEmoji = config.IconEmoji
		if t.IconEmoji == "" {
			t.IconEmoji = ":smile:"
		}
	}
}
