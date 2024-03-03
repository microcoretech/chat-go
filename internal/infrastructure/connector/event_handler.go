package connector

type EventHandler interface {
	HandleEvent(conn Connection, rawEvent Event) error
}
