package domain

type DeleteSidEvent struct {
	NetworkEvent
	key string
}

func NewDeleteSidEvent(key string) *DeleteSidEvent {
	return &DeleteSidEvent{
		key: key,
	}
}

func (deleteSidEvent *DeleteSidEvent) GetKey() string {
	return deleteSidEvent.key
}
