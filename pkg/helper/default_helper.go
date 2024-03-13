package helper

import (
	"github.com/hawkv6/hawkeye/pkg/logging"
	"github.com/sirupsen/logrus"
)

type DefaultHelper struct {
	log *logrus.Entry
}

func NewDefaultHelper() *DefaultHelper {
	return &DefaultHelper{
		log: logging.DefaultLogger.WithField("subsystem", subsystem),
	}
}

func (helper *DefaultHelper) GetLsNodeProperties() []string {
	helper.log.Debugln("Getting LsNode properties")
	return []string{"Key", "IgpRouterId", "Name"}
}

func (helper *DefaultHelper) GetLsLinkProperties() []string {
	helper.log.Debugln("Getting LsLink properties")
	return []string{
		"Key",
		"IgpRouterId",
		"RemoteIgpRouterId",
		"UnidirLinkDelay",
		"UnidirDelayVariation",
		"UnidirAvailableBW",
		"UnidirPacketLoss",
		"UnidirBWUtilization",
	}
}

func (helper *DefaultHelper) GetLsPrefixProperties() []string {
	helper.log.Debugln("Getting LsPrefix properties")
	return []string{"Key", "IgpRouterId", "Prefix", "PrefixLen"}
}

func (helper *DefaultHelper) GetLsSrv6SidsProperties() []string {
	helper.log.Debugln("Getting LsSrv6Sids properties")
	return []string{"Key", "IgpRouterId", "Srv6Sid"}
}
