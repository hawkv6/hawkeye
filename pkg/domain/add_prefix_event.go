package domain

type AddPrefixEvent struct {
	NetworkEvent
	Prefix
}

func NewAddPrefixEvent(prefix Prefix) *AddPrefixEvent {
	return &AddPrefixEvent{
		Prefix: prefix,
	}
}

func (addPrefixEvent *AddPrefixEvent) GetKey() string {
	return addPrefixEvent.Prefix.GetKey()
}
