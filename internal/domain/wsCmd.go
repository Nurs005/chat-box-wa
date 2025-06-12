package domain

type WebSocketCommand struct {
	Type    string `json:"type"`
	ChatJID string `json:"chat"`
	Text    string `json:"text"`
}
