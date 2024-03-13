package adapter

import (
	"context"

	"github.com/hawkv6/hawkeye/pkg/api"
	"github.com/hawkv6/hawkeye/pkg/domain"
	"github.com/jalapeno-api-gateway/jagw-go/jagw"
)

const Subsystem = "adapter"

type Adapter interface {
	ConvertNode(*jagw.LsNode) (domain.Node, error)
	ConvertNodeEvent(*jagw.LsNodeEvent) (domain.NetworkEvent, error)
	ConvertLink(*jagw.LsLink) (domain.Link, error)
	ConvertLinkEvent(*jagw.LsLinkEvent) (domain.NetworkEvent, error)
	ConvertPrefix(*jagw.LsPrefix) (domain.Prefix, error)
	ConvertPrefixEvent(*jagw.LsPrefixEvent) (domain.NetworkEvent, error)
	ConvertSid(*jagw.LsSrv6Sid) (domain.Sid, error)
	ConvertSidEvent(*jagw.LsSrv6SidEvent) (domain.NetworkEvent, error)
	ConvertPathRequest(*api.PathRequest, api.IntentController_GetIntentPathServer, context.Context) (domain.PathRequest, error)
	ConvertPathResult(domain.PathResult) (*api.PathResult, error)
}
