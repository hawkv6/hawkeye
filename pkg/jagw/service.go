package jagw

const Subsystem = "jagw"

type JagwService interface {
	Start() error
	Stop()
}
