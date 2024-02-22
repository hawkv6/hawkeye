package domain

import (
	"context"

	"github.com/hawkv6/hawkeye/pkg/api"
)

type PathRequest interface {
	GetIpv6SourceAddress() string
	GetIpv6DestinationAddress() string
	GetIntents() []Intent
	GetContext() context.Context
	GetStream() api.IntentController_GetIntentPathServer
}
