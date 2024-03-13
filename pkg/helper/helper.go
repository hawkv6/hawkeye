package helper

const subsystem = "helper"

type Helper interface {
	GetLsNodeProperties() []string
	GetLsLinkProperties() []string
	GetLsPrefixProperties() []string
	GetLsSrv6SidsProperties() []string
}
