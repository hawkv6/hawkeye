package domain

type DeleteLinkEvent struct {
	NetworkEvent
	key string
}

func NewDeleteLinkEvent(key string) *DeleteLinkEvent {
	return &DeleteLinkEvent{
		key: key,
	}
}

func (deleteLinkEvent *DeleteLinkEvent) GetKey() string {
	return deleteLinkEvent.key
}
