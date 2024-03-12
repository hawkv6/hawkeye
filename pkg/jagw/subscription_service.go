package jagw

type JagwSubscriptionService interface {
	JagwService
	SubscribeLsLinks() error
}
