package domain

type DeleteNodeEvent struct {
	NetworkEvent
	key string
}

func NewDeleteNodeEvent(key string) *DeleteNodeEvent {
	return &DeleteNodeEvent{
		key: key,
	}
}

func (deleteNodeEvent *DeleteNodeEvent) GetKey() string {
	return deleteNodeEvent.key
}
