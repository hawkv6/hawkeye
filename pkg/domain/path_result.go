package domain

type PathResult interface {
	PathRequest
	GetIpv6SidAddresses() []string
}
