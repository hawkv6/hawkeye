package domain

type UpdateNodeEvent struct {
	NetworkEvent
	Node
}

func NewUpdateNodeEvent(node Node) *UpdateNodeEvent {
	return &UpdateNodeEvent{
		Node: node,
	}
}

func (updateNodeEvent *UpdateNodeEvent) GetKey() string {
	return updateNodeEvent.Node.GetKey()
}
