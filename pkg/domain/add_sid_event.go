package domain

type AddSidEvent struct {
	NetworkEvent
	Sid
}

func NewAddSidEvent(sid Sid) *AddSidEvent {
	return &AddSidEvent{
		Sid: sid,
	}
}

func (addSidEvent *AddSidEvent) GetKey() string {
	return addSidEvent.Sid.GetKey()
}
