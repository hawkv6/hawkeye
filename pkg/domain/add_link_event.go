package domain

type AddLinkEvent struct {
	NetworkEvent
	Link
}

func NewAddLinkEvent(link Link) *AddLinkEvent {
	return &AddLinkEvent{
		Link: link,
	}
}

func (addLinkEvent *AddLinkEvent) GetKey() string {
	return addLinkEvent.Link.GetKey()
}
