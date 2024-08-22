package controller

const Subsystem = "controller"

type Controller interface {
	Start()
	Stop()
}
