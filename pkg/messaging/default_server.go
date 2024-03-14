package messaging

import (
	"context"
	"fmt"
	"io"
	"net"

	"github.com/hawkv6/hawkeye/pkg/adapter"
	"github.com/hawkv6/hawkeye/pkg/api"
	"github.com/hawkv6/hawkeye/pkg/config"
	"github.com/hawkv6/hawkeye/pkg/domain"
	"github.com/hawkv6/hawkeye/pkg/logging"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
)

type DefaultMessagingServer struct {
	api.UnimplementedIntentControllerServer
	log             *logrus.Entry
	adapter         adapter.Adapter
	grpcPort        uint16
	pathRequestChan chan domain.PathRequest
	pathResultChan  chan domain.PathResult
}

func NewDefaultMessagingServer(adapter adapter.Adapter, config config.Config, messagingChannels MessagingChannels) *DefaultMessagingServer {
	return &DefaultMessagingServer{
		log:             logging.DefaultLogger.WithField("subsystem", Subsystem),
		adapter:         adapter,
		grpcPort:        config.GetGrpcPort(),
		pathRequestChan: messagingChannels.GetPathRequestChan(),
		pathResultChan:  messagingChannels.GetPathResponseChan(),
	}
}

func (server *DefaultMessagingServer) Start() error {
	listenAddress := fmt.Sprintf(":%d", server.grpcPort)
	list, err := net.Listen("tcp", listenAddress)
	if err != nil {
		return fmt.Errorf("Failed to listen: %v", err)
	}
	server.log.Infoln("Listening on " + listenAddress)

	grpcServer := grpc.NewServer()
	api.RegisterIntentControllerServer(grpcServer, server)
	if err := grpcServer.Serve(list); err != nil {
		return fmt.Errorf("Failed to serve: %v", err)
	}
	return nil
}

func (server *DefaultMessagingServer) GetIntentPath(stream api.IntentController_GetIntentPathServer) error {
	ctx := stream.Context()
	peerInfo, ok := peer.FromContext(ctx)
	if ok {
		server.log.Debugln("Received Stream from: ", peerInfo.Addr)
	}
	for {
		apiRequest, err := stream.Recv()
		if err != nil {
			ctx.Done()
			if err == io.EOF {
				server.log.Debugf("Stream has with %s ended", peerInfo.Addr)
				return nil
			} else {
				server.log.Errorln("Error receiving message: ", err)
				return err
			}
		}
		server.log.Debugln("Received request: ", apiRequest)

		pathRequest, err := server.adapter.ConvertPathRequest(apiRequest, stream, ctx)
		if err != nil {
			return err
		}
		server.pathRequestChan <- pathRequest
		go server.GetIntentPathResponse(stream, ctx)
	}
}

func (server *DefaultMessagingServer) GetIntentPathResponse(stream api.IntentController_GetIntentPathServer, ctx context.Context) {
	for {
		select {
		case pathResult := <-server.pathResultChan:
			result, err := server.adapter.ConvertPathResult(pathResult)
			if err != nil {
				server.log.Errorln("Error converting PathResult: ", err)
				return
			}
			if err := stream.Send(result); err != nil {
				server.log.Errorln("Error sending message: ", err)
				return
			}
		case <-ctx.Done():
			return
		}
	}
}
