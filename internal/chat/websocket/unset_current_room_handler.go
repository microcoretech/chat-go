package websocket

func (e *EventHandler) unsetCurrentChatHandler(conn Connection, rawData []byte) error {
	conn.SetSubscribedChats(nil)
	return nil
}
