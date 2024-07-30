package domain

type AddNodeEvent struct {
	NetworkEvent
	Node
}

func NewAddNodeEvent(node Node) *AddNodeEvent {
	return &AddNodeEvent{
		Node: node,
	}
}

func (addNodeEvent *AddNodeEvent) GetKey() string {
	return addNodeEvent.Node.GetKey()
}
