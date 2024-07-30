package jagw

const Subsystem = "jagw"

type JagwService interface {
	Init() error
	Start() error
	Stop()
}
