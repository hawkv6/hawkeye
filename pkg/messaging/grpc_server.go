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

type GrpcMessagingServer struct {
	api.UnimplementedIntentControllerServer
	log             *logrus.Entry
	adapter         adapter.Adapter
	grpcPort        uint16
	pathRequestChan chan domain.PathRequest
	pathResultChan  chan domain.PathResult
	errorChan       chan error
	internalChan    chan error
	stopChan        chan struct{}
}

func NewGrpcMessagingServer(adapter adapter.Adapter, config config.Config, messagingChannels MessagingChannels) *GrpcMessagingServer {
	return &GrpcMessagingServer{
		log:             logging.DefaultLogger.WithField("subsystem", Subsystem),
		adapter:         adapter,
		grpcPort:        config.GetGrpcPort(),
		pathRequestChan: messagingChannels.GetPathRequestChan(),
		pathResultChan:  messagingChannels.GetPathResponseChan(),
		errorChan:       messagingChannels.GetErrorChan(),
		internalChan:    make(chan error),
		stopChan:        make(chan struct{}),
	}
}

func (server *GrpcMessagingServer) Start() error {
	listenAddress := fmt.Sprintf(":%d", server.grpcPort)
	list, err := net.Listen("tcp", listenAddress)
	if err != nil {
		return fmt.Errorf("Failed to listen: %v", err)
	}
	server.log.Infoln("Listening on " + listenAddress)

	grpcServer := grpc.NewServer()
	api.RegisterIntentControllerServer(grpcServer, server)

	go func() {
		if err := grpcServer.Serve(list); err != nil {
			server.log.Fatalf("Error starting gRPC server %v", err)
		}
	}()

	<-server.stopChan
	grpcServer.GracefulStop()
	return nil
}

func (server *GrpcMessagingServer) handleIncomingPathRequests(stream api.IntentController_GetIntentPathServer, peerInfo *peer.Peer, ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			server.log.Debugln("Context cancelled, stopping receiving the stream")
			return
		default:
			apiRequest, err := stream.Recv()
			if err != nil {
				if err == io.EOF {
					server.log.Debugf("Stream %s ended", peerInfo.Addr)
					return
				} else {
					server.log.Errorln("Error receiving message: ", err)
					server.internalChan <- err
					return
				}
			}
			server.log.Debugln("Received request: ", apiRequest)
			pathRequest, err := server.adapter.ConvertPathRequest(apiRequest, stream, ctx)
			if err != nil {
				server.log.Errorln("Error converting PathRequest: ", err)
				server.internalChan <- err
				return
			}
			server.pathRequestChan <- pathRequest
			go server.handleIntentPathResponse(stream, ctx)
		}
	}
}

func (server *GrpcMessagingServer) GetIntentPath(stream api.IntentController_GetIntentPathServer) error {
	ctx := stream.Context()
	peerInfo, ok := peer.FromContext(ctx)
	if ok {
		server.log.Debugln("Received Stream from: ", peerInfo.Addr)
	}
	go server.handleIncomingPathRequests(stream, peerInfo, ctx)
	select {
	case <-ctx.Done():
		return nil
	case err := <-server.internalChan:
		return err
	}
}

func (server *GrpcMessagingServer) handleIntentPathResponse(stream api.IntentController_GetIntentPathServer, ctx context.Context) {
	for {
		select {
		case pathResult := <-server.pathResultChan:
			result, err := server.adapter.ConvertPathResult(pathResult)
			if err != nil {
				server.log.Errorln("Error converting PathResult: ", err)
				server.internalChan <- err
				return
			}
			if err := stream.Send(result); err != nil {
				server.log.Errorln("Error sending message: ", err)
				server.internalChan <- err
				return
			}
		case err := <-server.errorChan:
			server.internalChan <- err
		case <-ctx.Done():
			server.log.Debugln("Context cancelled, stopping GetIntentPathResponse")
			return
		}
	}
}

func (server *GrpcMessagingServer) Stop() {
	server.log.Infoln("Stopping the gRPC server")
	close(server.stopChan)
}