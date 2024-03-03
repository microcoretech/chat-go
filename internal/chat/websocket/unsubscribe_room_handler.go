package websocket

func (e *EventHandler) unsubscribeRoomHandler(conn Connection, rawData []byte) error {
	conn.SetCurrentChat(nil)
	return nil
}
