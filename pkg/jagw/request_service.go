package jagw

type JagwRequestService interface {
	JagwService
	GetLsLinks() error
	GetLsPrefixes() error
}
