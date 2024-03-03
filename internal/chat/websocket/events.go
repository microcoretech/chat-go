package websocket

const (
	SubscribeChatsEventType       = 1
	UnsubscribeChatsEventType     = 2
	SetCurrentChatEventType       = 3
	UnsetCurrentChatEventType     = 4
	CreateMessageEventType        = 5
	EditMessageEventType          = 6
	DeleteMessageEventType        = 7
	UpdateMessagesStatusEventType = 8
)

type EditMessageEventData struct {
	MessageID string `json:"messageId"`
	Text      string `json:"text"`
}

type DeleteMessageEventData struct {
	MessageID string `json:"messageId"`
}
