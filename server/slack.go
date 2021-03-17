package server

type SlackResponse struct {
	Text         string `json:"text"`
	ResponseType string `json:"response_type"`
}

const (
	Ephemeral = "ephemeral"
	InChannel = "in_channel"
)
