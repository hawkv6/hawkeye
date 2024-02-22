package domain

type Prefix interface {
	GetKey() string
	GetIgpRouterId() string
	GetPrefix() string
	GetPrefixLength() uint8
}
