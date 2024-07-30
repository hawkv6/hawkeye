package domain

type UpdateLinkEvent struct {
	NetworkEvent
	Link
}

func NewUpdateLinkEvent(link Link) *UpdateLinkEvent {
	return &UpdateLinkEvent{
		Link: link,
	}
}

func (updateLinkEvent *UpdateLinkEvent) GetKey() string {
	return updateLinkEvent.Link.GetKey()
}
