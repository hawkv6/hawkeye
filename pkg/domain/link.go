package domain

type Link interface {
	GetKey() string
	GetIgpRouterId() string
	GetRemoteIgpRouterId() string
	GetUnidirLinkDelay() float64
}
