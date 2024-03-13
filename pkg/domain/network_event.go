package domain

type NetworkEvent interface {
	GetKey() string
}

// type LinkEvent struct {
// 	Event
// 	Link *Link
// }

// type PrefixEvent struct {
// 	Event
// 	Prefix *Prefix
// }

// type SidEvent struct {
// 	Event
// 	Sid *Sid
// }
