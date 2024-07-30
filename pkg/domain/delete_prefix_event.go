package domain

type DeletePrefixEvent struct {
	NetworkEvent
	key string
}

func NewDeletePrefixEvent(key string) *DeletePrefixEvent {
	return &DeletePrefixEvent{
		key: key,
	}
}

func (deletePrefixEvent *DeletePrefixEvent) GetKey() string {
	return deletePrefixEvent.key
}
