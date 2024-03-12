package jagw

type JagwRequestService interface {
	JagwService
	GetLsNodes() error
	GetLsLinks() error
	GetLsPrefixes() error
	GetSrv6Sids() error
}
