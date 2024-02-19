package jagw

const Subsystem = "jagw"

type JagwService interface {
	Init() error
}

type JagwRequestService interface {
	JagwService
	GetLsNodeEdge() error
}

type JagwSubscriptionService interface {
	JagwService
}
