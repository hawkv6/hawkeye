package adapter

import (
	"context"

	"github.com/hawkv6/hawkeye/pkg/api"
	"github.com/hawkv6/hawkeye/pkg/domain"
	"github.com/jalapeno-api-gateway/jagw-go/jagw"
)

const Subsystem = "adapter"

type Adapter interface {
	ConvertLink(*jagw.LsLink) (domain.Link, error)
	ConvertPrefix(*jagw.LsPrefix) (domain.Prefix, error)
	ConvertSid(*jagw.LsSrv6Sid) (domain.Sid, error)
	ConvertPathRequest(*api.PathRequest, api.IntentController_GetIntentPathServer, context.Context) (domain.PathRequest, error)
	ConvertPathResult(domain.PathResult) (*api.PathResult, error)
}
