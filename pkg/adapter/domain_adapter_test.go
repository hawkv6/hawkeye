package adapter

import (
	"context"
	"reflect"
	"testing"

	"github.com/hawkv6/hawkeye/pkg/api"
	"github.com/hawkv6/hawkeye/pkg/domain"
	"github.com/hawkv6/hawkeye/pkg/graph"
	"github.com/hawkv6/hawkeye/pkg/logging"
	"github.com/jalapeno-api-gateway/jagw-go/jagw"
	"github.com/sirupsen/logrus"
	"go.uber.org/mock/gomock"
	"google.golang.org/protobuf/proto"
)

func TestNewDomainAdapter(t *testing.T) {
	tests := []struct {
		name    string
		wantNil bool
	}{
		{
			name:    "Create new DomainAdapter",
			wantNil: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewDomainAdapter()
			if (got == nil) != tt.wantNil {
				t.Errorf("NewDomainAdapter() = %v, want non-nil %v", got, !tt.wantNil)
			}
		})
	}
}

func setUpJagwNode(key string, igpRouterId string, name string, srAlgorithm []uint32) *jagw.LsNode {
	lsNode := &jagw.LsNode{}
	if key != "" {
		lsNode.Key = proto.String(key)
	}
	if igpRouterId != "" {
		lsNode.IgpRouterId = proto.String(igpRouterId)
	}
	if name != "" {
		lsNode.Name = proto.String(name)
	}
	if srAlgorithm != nil {
		lsNode.SrAlgorithm = srAlgorithm
	}
	return lsNode
}

func setUpDomainNode(key string, igpRouterId string, name string, srAlgorithm []uint32) *domain.DomainNode {
	node, _ := domain.NewDomainNode(&key, &igpRouterId, &name, srAlgorithm)
	return node
}

func isNilInterface(value interface{}) bool {
	return value == nil || reflect.ValueOf(value).IsNil()
}

func TestDomainAdapter_ConvertNode(t *testing.T) {
	type fields struct {
		log *logrus.Entry
	}
	type args struct {
		lsNode *jagw.LsNode
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    domain.Node
		wantErr bool
	}{
		{
			name: "Convert LsNode to Node successfully",
			fields: fields{
				log: logging.DefaultLogger.WithField("subsystem", Subsystem),
			},
			args: args{
				lsNode: setUpJagwNode("key", "igpRouterId", "name", []uint32{1, 2, 3}),
			},
			want:    setUpDomainNode("key", "igpRouterId", "name", []uint32{1, 2, 3}),
			wantErr: false,
		},
		{
			name: "Convert LsNode to Node no key",
			fields: fields{
				log: logging.DefaultLogger.WithField("subsystem", Subsystem),
			},
			args: args{
				lsNode: setUpJagwNode("", "igpRouterId", "name", []uint32{1, 2, 3}),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Convert LsNode to Node no igp router id",
			fields: fields{
				log: logging.DefaultLogger.WithField("subsystem", Subsystem),
			},
			args: args{
				lsNode: setUpJagwNode("key", "", "name", []uint32{1, 2, 3}),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Convert LsNode to Node no name",
			fields: fields{
				log: logging.DefaultLogger.WithField("subsystem", Subsystem),
			},
			args: args{
				lsNode: setUpJagwNode("key", "igpRouterId", "", []uint32{1, 2, 3}),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Convert LsNode to Node no name",
			fields: fields{
				log: logging.DefaultLogger.WithField("subsystem", Subsystem),
			},
			args: args{
				lsNode: setUpJagwNode("key", "igpRouterId", "name", nil),
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			adapter := &DomainAdapter{
				log: tt.fields.log,
			}
			got, err := adapter.ConvertNode(tt.args.lsNode)
			if (err != nil) != tt.wantErr {
				t.Errorf("DomainAdapter.ConvertNode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if isNilInterface(got) && tt.want != nil {
				t.Errorf("DomainAdapter.ConvertNode() = %v, want %v", got, tt.want)
			}
			if !isNilInterface(got) && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DomainAdapter.ConvertNode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDomainAdapter_ConvertNodeEvent(t *testing.T) {
	type fields struct {
		log *logrus.Entry
	}
	tests := []struct {
		fields      fields
		name        string
		lsNodeEvent *jagw.LsNodeEvent
		wantEvent   domain.NetworkEvent
		wantErr     bool
	}{
		{
			name:        "nil LsNodeEvent",
			lsNodeEvent: nil,
			wantErr:     true,
		},
		{
			name:        "delete action success",
			lsNodeEvent: &jagw.LsNodeEvent{Action: proto.String("del"), Key: proto.String("key")},
			wantEvent:   domain.NewDeleteNodeEvent("key"),
			wantErr:     false,
		},
		{
			name:        "add action with ConvertNode error",
			lsNodeEvent: &jagw.LsNodeEvent{Key: proto.String("key"), Action: proto.String("add"), LsNode: setUpJagwNode("", "igpRouterId", "name", []uint32{1, 2, 3})},
			wantErr:     true,
		},
		{
			name:        "add action success",
			lsNodeEvent: &jagw.LsNodeEvent{Key: proto.String("key"), Action: proto.String("add"), LsNode: setUpJagwNode("key", "igpRouterId", "name", []uint32{1, 2, 3})},
			wantEvent:   domain.NewAddNodeEvent(setUpDomainNode("key", "igpRouterId", "name", []uint32{1, 2, 3})),
			wantErr:     false,
		},
		{
			name:        "update action with ConvertNode error",
			lsNodeEvent: &jagw.LsNodeEvent{Key: proto.String("key"), Action: proto.String("update"), LsNode: setUpJagwNode("", "igpRouterId", "name", []uint32{1, 2, 3})},
			wantErr:     true,
		},
		{
			name:        "update action success",
			lsNodeEvent: &jagw.LsNodeEvent{Key: proto.String("key"), Action: proto.String("update"), LsNode: setUpJagwNode("key", "igpRouterId", "name", []uint32{1, 2, 3})},
			wantEvent:   domain.NewUpdateNodeEvent(setUpDomainNode("key", "igpRouterId", "name", []uint32{1, 2, 3})),
			wantErr:     false,
		},
		{
			name:        "LsNodeEvent with nil action",
			lsNodeEvent: &jagw.LsNodeEvent{Key: proto.String("key"), LsNode: setUpJagwNode("key", "igpRouterId", "name", []uint32{1, 2, 3})},
			wantEvent:   nil,
			wantErr:     true,
		},
		{
			name:        "Uknown action",
			lsNodeEvent: &jagw.LsNodeEvent{Key: proto.String("key"), Action: proto.String("Uknown"), LsNode: setUpJagwNode("key", "igpRouterId", "name", []uint32{1, 2, 3})},
			wantEvent:   nil,
			wantErr:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			adapter := &DomainAdapter{
				log: tt.fields.log,
			}
			gotEvent, err := adapter.ConvertNodeEvent(tt.lsNodeEvent)
			if (err != nil) != tt.wantErr {
				t.Errorf("ConvertNodeEvent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotEvent, tt.wantEvent) {
				t.Errorf("ConvertNodeEvent() gotEvent = %v, want %v", gotEvent, tt.wantEvent)
			}
		})
	}
}

func setUpJagwLink(key string, igpRouterId string, remoteIgpRouterId string, igpMetric uint32, unidirLinkDelay uint32, unidirDelayVariation uint32, maxLinkBwKbps uint64, unidirAvailableBw uint32, unidirBwUtilization uint32, unidirPacketLossPercentage, normalizedUnidirLinkDelay, normalizedUnidirDelayVariation, normalizedUnidirPacketLoss float64) *jagw.LsLink {
	lsLink := &jagw.LsLink{}
	if key != "" {
		lsLink.Key = proto.String(key)
	}
	if igpRouterId != "" {
		lsLink.IgpRouterId = proto.String(igpRouterId)
	}
	if remoteIgpRouterId != "" {
		lsLink.RemoteIgpRouterId = proto.String(remoteIgpRouterId)
	}
	lsLink.IgpMetric = proto.Uint32(igpMetric)
	lsLink.UnidirLinkDelay = proto.Uint32(unidirLinkDelay)
	lsLink.UnidirDelayVariation = proto.Uint32(unidirDelayVariation)
	lsLink.MaxLinkBwKbps = proto.Uint64(maxLinkBwKbps)
	lsLink.UnidirAvailableBw = proto.Uint32(unidirAvailableBw)
	lsLink.UnidirBwUtilization = proto.Uint32(unidirBwUtilization)
	lsLink.UnidirPacketLossPercentage = proto.Float64(unidirPacketLossPercentage)
	lsLink.NormalizedUnidirLinkDelay = proto.Float64(normalizedUnidirLinkDelay)
	lsLink.NormalizedUnidirDelayVariation = proto.Float64(normalizedUnidirDelayVariation)
	lsLink.NormalizedUnidirPacketLoss = proto.Float64(normalizedUnidirPacketLoss)
	return lsLink
}

func setUpDomainLink(key string, igpRouterId string, remoteIgpRouterId string, igpMetric uint32, unidirLinkDelay uint32, unidirDelayVariation uint32, maxLinkBwKbps uint64, unidirAvailableBw uint32, unidirBwUtilization uint32, unidirPacketLossPercentage, normalizedUnidirLinkDelay, normalizedUnidirDelayVariation, normalizedUnidirPacketLoss float64) *domain.DomainLink {
	link, _ := domain.NewDomainLink(&key, &igpRouterId, &remoteIgpRouterId, &igpMetric, &unidirLinkDelay, &unidirDelayVariation, &maxLinkBwKbps, &unidirAvailableBw, &unidirBwUtilization, &unidirPacketLossPercentage, &normalizedUnidirLinkDelay, &normalizedUnidirDelayVariation, &normalizedUnidirPacketLoss)
	return link
}

func TestDomainAdapter_ConvertLink(t *testing.T) {
	type fields struct {
		log *logrus.Entry
	}
	type args struct {
		lsLink *jagw.LsLink
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    domain.Link
		wantErr bool
	}{
		{
			name: "Convert LsLink to Link successfully",
			fields: fields{
				log: logging.DefaultLogger.WithField("subsystem", Subsystem),
			},
			args: args{
				lsLink: setUpJagwLink("key", "igpRouterId", "remoteIgpRouterId", 1, 2, 3, 1000, 100, 50, 0.5, 1.0, 1.0, 1.0),
			},
			want:    setUpDomainLink("key", "igpRouterId", "remoteIgpRouterId", 1, 2, 3, 1000, 100, 50, 0.5, 1.0, 1.0, 1.0),
			wantErr: false,
		},
		{
			name: "Convert LsLink to Link Error:Field validation for 'MaxLinkBWKbps' failed on the 'min",
			fields: fields{
				log: logging.DefaultLogger.WithField("subsystem", Subsystem),
			},
			args: args{
				lsLink: setUpJagwLink("key", "igpRouterId", "remoteIgpRouterId", 1, 2, 3, 0, 100, 50, 0.5, 1.0, 1.0, 1.0),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Convert LsLink to Link Error:Field validation for 'UnidirPacketLoss' failed on the 'min",
			fields: fields{
				log: logging.DefaultLogger.WithField("subsystem", Subsystem),
			},
			args: args{
				lsLink: setUpJagwLink("key", "igpRouterId", "remoteIgpRouterId", 1, 2, 3, 1000, 100, 50, -0.5, 1.0, 1.0, 1.0),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Convert LsLink to Link Error:Field validation for 'UnidirPacketLoss' failed on the 'max",
			fields: fields{
				log: logging.DefaultLogger.WithField("subsystem", Subsystem),
			},
			args: args{
				lsLink: setUpJagwLink("key", "igpRouterId", "remoteIgpRouterId", 1, 2, 3, 1000, 100, 50, 100.1, 1.0, 1.0, 1.0),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Convert LsLink to Link Error:Field validation for 'NormalizedUnidirLinkDelay' failed on the 'min'",
			fields: fields{
				log: logging.DefaultLogger.WithField("subsystem", Subsystem),
			},
			args: args{
				lsLink: setUpJagwLink("key", "igpRouterId", "remoteIgpRouterId", 1, 2, 3, 1000, 100, 50, 0.5, -1.0, 1.0, 1.0),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Convert LsLink to Link Error:Field validation for 'NormalizedUnidirLinkDelay' failed on the 'max'",
			fields: fields{
				log: logging.DefaultLogger.WithField("subsystem", Subsystem),
			},
			args: args{
				lsLink: setUpJagwLink("key", "igpRouterId", "remoteIgpRouterId", 1, 2, 3, 1000, 100, 50, 0.5, 2.0, 1.0, 1.0),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Convert LsLink to Link Error:Field validation for 'NormalizedUnidirDelayVariaton' failed on the 'min'",
			fields: fields{
				log: logging.DefaultLogger.WithField("subsystem", Subsystem),
			},
			args: args{
				lsLink: setUpJagwLink("key", "igpRouterId", "remoteIgpRouterId", 1, 2, 3, 1000, 100, 50, 0.5, 1.0, -1.0, 1.0),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Convert LsLink to Link Error:Field validation for 'NormalizedUnidirLinkDelayVariation' failed on the 'max'",
			fields: fields{
				log: logging.DefaultLogger.WithField("subsystem", Subsystem),
			},
			args: args{
				lsLink: setUpJagwLink("key", "igpRouterId", "remoteIgpRouterId", 1, 2, 3, 1000, 100, 50, 0.5, 1.0, 2.0, 1.0),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Convert LsLink to Link Error:Field validation for 'NormalizedUnidirPacketLoss' failed on the 'min'",
			fields: fields{
				log: logging.DefaultLogger.WithField("subsystem", Subsystem),
			},
			args: args{
				lsLink: setUpJagwLink("key", "igpRouterId", "remoteIgpRouterId", 1, 2, 3, 1000, 100, 50, 0.5, 1.0, 1.0, -1.0),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Convert LsLink to Link Error:Field validation for 'NormalizedUnidirLinkPacketLoss' failed on the 'max'",
			fields: fields{
				log: logging.DefaultLogger.WithField("subsystem", Subsystem),
			},
			args: args{
				lsLink: setUpJagwLink("key", "igpRouterId", "remoteIgpRouterId", 1, 2, 3, 1000, 100, 50, 0.5, 1.0, 1.0, 2.0),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Convert LsLink to Link no key",
			fields: fields{
				log: logging.DefaultLogger.WithField("subsystem", Subsystem),
			},
			args: args{
				lsLink: setUpJagwLink("", "igpRouterId", "remoteIgpRouterId", 1, 2, 3, 1000, 100, 50, 0.5, 1.0, 1.5, 2.0),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Convert LsLink to Link no igp router id",
			fields: fields{
				log: logging.DefaultLogger.WithField("subsystem", Subsystem),
			},
			args: args{
				lsLink: setUpJagwLink("key", "", "remoteIgpRouterId", 1, 2, 3, 1000, 100, 50, 0.5, 1.0, 1.5, 2.0),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Convert LsLink to Link no remote igp router id",
			fields: fields{
				log: logging.DefaultLogger.WithField("subsystem", Subsystem),
			},
			args: args{
				lsLink: setUpJagwLink("key", "igpRouterId", "", 1, 2, 3, 1000, 100, 50, 0.5, 1.0, 1.5, 2.0),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Convert LsLink to Link with nil igp metric",
			fields: fields{
				log: logging.DefaultLogger.WithField("subsystem", Subsystem),
			},
			args: args{
				lsLink: setUpJagwLink("key", "igpRouterId", "remoteIgpRouterId", 0, 2, 3, 1000, 100, 50, 0.5, 1.0, 1.5, 2.0),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Convert LsLink to Link with nil unidir link delay",
			fields: fields{
				log: logging.DefaultLogger.WithField("subsystem", Subsystem),
			},
			args: args{
				lsLink: setUpJagwLink("key", "igpRouterId", "remoteIgpRouterId", 1, 0, 3, 1000, 100, 50, 0.5, 1.0, 1.5, 2.0),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Convert LsLink to Link with nil unidir delay variation",
			fields: fields{
				log: logging.DefaultLogger.WithField("subsystem", Subsystem),
			},
			args: args{
				lsLink: setUpJagwLink("key", "igpRouterId", "remoteIgpRouterId", 1, 2, 0, 1000, 100, 50, 0.5, 1.0, 1.5, 2.0),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Convert LsLink to Link with nil max link bw kbps",
			fields: fields{
				log: logging.DefaultLogger.WithField("subsystem", Subsystem),
			},
			args: args{
				lsLink: setUpJagwLink("key", "igpRouterId", "remoteIgpRouterId", 1, 2, 3, 0, 100, 50, 0.5, 1.0, 1.5, 2.0),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Convert LsLink to Link with nil unidir available bw",
			fields: fields{
				log: logging.DefaultLogger.WithField("subsystem", Subsystem),
			},
			args: args{
				lsLink: setUpJagwLink("key", "igpRouterId", "remoteIgpRouterId", 1, 2, 3, 1000, 0, 50, 0.5, 1.0, 1.5, 2.0),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Convert LsLink to Link with nil unidir bw utilization",
			fields: fields{
				log: logging.DefaultLogger.WithField("subsystem", Subsystem),
			},
			args: args{
				lsLink: setUpJagwLink("key", "igpRouterId", "remoteIgpRouterId", 1, 2, 3, 1000, 100, 0, 0.5, 1.0, 1.5, 2.0),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Convert LsLink to Link with nil unidir packet loss percentage",
			fields: fields{
				log: logging.DefaultLogger.WithField("subsystem", Subsystem),
			},
			args: args{
				lsLink: setUpJagwLink("key", "igpRouterId", "remoteIgpRouterId", 1, 2, 3, 1000, 100, 50, 0.0, 1.0, 1.5, 2.0),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Convert LsLink to Link with nil normalized unidir link delay",
			fields: fields{
				log: logging.DefaultLogger.WithField("subsystem", Subsystem),
			},
			args: args{
				lsLink: setUpJagwLink("key", "igpRouterId", "remoteIgpRouterId", 1, 2, 3, 1000, 100, 50, 0.5, 0.0, 1.5, 2.0),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Convert LsLink to Link with nil normalized unidir delay variation",
			fields: fields{
				log: logging.DefaultLogger.WithField("subsystem", Subsystem),
			},
			args: args{
				lsLink: setUpJagwLink("key", "igpRouterId", "remoteIgpRouterId", 1, 2, 3, 1000, 100, 50, 0.5, 1.0, 0.0, 2.0),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Convert LsLink to Link with nil normalized unidir packet loss",
			fields: fields{
				log: logging.DefaultLogger.WithField("subsystem", Subsystem),
			},
			args: args{
				lsLink: setUpJagwLink("key", "igpRouterId", "remoteIgpRouterId", 1, 2, 3, 1000, 100, 50, 0.5, 1.0, 1.5, 0.0),
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			adapter := &DomainAdapter{
				log: tt.fields.log,
			}
			got, err := adapter.ConvertLink(tt.args.lsLink)
			if (err != nil) != tt.wantErr {
				t.Errorf("DomainAdapter.ConvertLink() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if isNilInterface(got) && tt.want != nil {
				t.Errorf("DomainAdapter.ConvertLink() = %v, want %v", got, tt.want)
			}
			if !isNilInterface(got) && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DomainAdapter.ConvertLink() = %v, want %v", got, tt.want)
			}
		})
	}
}
func TestDomainAdapter_ConvertLinkEvent(t *testing.T) {
	type fields struct {
		log *logrus.Entry
	}
	tests := []struct {
		fields      fields
		name        string
		lsLinkEvent *jagw.LsLinkEvent
		want        domain.NetworkEvent
		wantErr     bool
	}{
		{
			name:        "nil LsLinkEvent",
			lsLinkEvent: nil,
			wantErr:     true,
		},
		{
			name:        "delete action success",
			lsLinkEvent: &jagw.LsLinkEvent{Action: proto.String("del"), Key: proto.String("key")},
			want:        domain.NewDeleteLinkEvent("key"),
			wantErr:     false,
		},
		{
			name:        "add action with ConvertLink error",
			lsLinkEvent: &jagw.LsLinkEvent{Key: proto.String("key"), Action: proto.String("add"), LsLink: setUpJagwLink("", "igpRouterId", "remoteIgpRouterId", 1, 2, 3, 1000, 100, 50, 0.5, 1.0, 1.0, 1.0)},
			wantErr:     true,
		},
		{
			name:        "add action success",
			lsLinkEvent: &jagw.LsLinkEvent{Key: proto.String("key"), Action: proto.String("add"), LsLink: setUpJagwLink("key", "igpRouterId", "remoteIgpRouterId", 1, 2, 3, 1000, 100, 50, 0.5, 1.0, 1.0, 1.0)},
			want:        domain.NewAddLinkEvent(setUpDomainLink("key", "igpRouterId", "remoteIgpRouterId", 1, 2, 3, 1000, 100, 50, 0.5, 1.0, 1.0, 1.0)),
			wantErr:     false,
		},
		{
			name:        "update action with ConvertLink error",
			lsLinkEvent: &jagw.LsLinkEvent{Key: proto.String("key"), Action: proto.String("update"), LsLink: setUpJagwLink("", "igpRouterId", "remoteIgpRouterId", 1, 2, 3, 1000, 100, 50, 0.5, 1.0, 1.0, 1.0)},
			wantErr:     true,
		},
		{
			name:        "update action success",
			lsLinkEvent: &jagw.LsLinkEvent{Key: proto.String("key"), Action: proto.String("update"), LsLink: setUpJagwLink("key", "igpRouterId", "remoteIgpRouterId", 1, 2, 3, 1000, 100, 50, 0.5, 1.0, 1.0, 1.0)},
			want:        domain.NewUpdateLinkEvent(setUpDomainLink("key", "igpRouterId", "remoteIgpRouterId", 1, 2, 3, 1000, 100, 50, 0.5, 1.0, 1.0, 1.0)),
			wantErr:     false,
		},
		{
			name:        "LsLinkEvent with nil action",
			lsLinkEvent: &jagw.LsLinkEvent{Key: proto.String("key"), LsLink: setUpJagwLink("key", "igpRouterId", "remoteIgpRouterId", 1, 2, 3, 1000, 100, 50, 0.5, 1.0, 1.0, 1.0)},
			want:        nil,
			wantErr:     true,
		},
		{
			name:        "Uknown action",
			lsLinkEvent: &jagw.LsLinkEvent{Key: proto.String("key"), Action: proto.String("Uknown"), LsLink: setUpJagwLink("key", "igpRouterId", "remoteIgpRouterId", 1, 2, 3, 1000, 100, 50, 0.5, 1.0, 1.0, 1.0)},
			want:        nil,
			wantErr:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			adapter := &DomainAdapter{
				log: tt.fields.log,
			}
			got, err := adapter.ConvertLinkEvent(tt.lsLinkEvent)
			if (err != nil) != tt.wantErr {
				t.Errorf("DomainAdapter.ConvertLinkEvent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if isNilInterface(got) && tt.want != nil {
				t.Errorf("DomainAdapter.ConvertLinkEvent() = %v, want %v", got, tt.want)
			}
			if !isNilInterface(got) && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DomainAdapter.ConvertLinkEvent() = %v, want %v", got, tt.want)
			}
		})
	}
}

func setUpJagwPrefix(key string, igpRouterId string, prefix string, prefixLength int32) *jagw.LsPrefix {
	lsPrefix := &jagw.LsPrefix{}
	if key != "" {
		lsPrefix.Key = proto.String(key)
	}
	if igpRouterId != "" {
		lsPrefix.IgpRouterId = proto.String(igpRouterId)
	}
	if prefix != "" {
		lsPrefix.Prefix = proto.String(prefix)
	}
	if prefixLength != 0 {
		lsPrefix.PrefixLen = proto.Int32(prefixLength)
	}
	return lsPrefix
}

func setUpDomainPrefix(key string, igpRouterId string, prefixValue string, prefixLength int32) *domain.DomainPrefix {
	prefix, _ := domain.NewDomainPrefix(&key, &igpRouterId, &prefixValue, &prefixLength)
	return prefix
}

func TestDomainAdapter_ConvertPrefix(t *testing.T) {
	type fields struct {
		log *logrus.Entry
	}
	type args struct {
		lsPrefix *jagw.LsPrefix
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    domain.Prefix
		wantErr bool
	}{
		{
			name: "Convert LsPrefix to Prefix successfully",
			fields: fields{
				log: logging.DefaultLogger.WithField("subsystem", Subsystem),
			},
			args: args{
				lsPrefix: setUpJagwPrefix("key", "igpRouterId", "prefix", 24),
			},
			want:    setUpDomainPrefix("key", "igpRouterId", "prefix", 24),
			wantErr: false,
		},
		{
			name: "Convert LsPrefix to Prefix no key",
			fields: fields{
				log: logging.DefaultLogger.WithField("subsystem", Subsystem),
			},
			args: args{
				lsPrefix: setUpJagwPrefix("", "igpRouterId", "prefix", 24),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Convert LsPrefix to Prefix no igp router id",
			fields: fields{
				log: logging.DefaultLogger.WithField("subsystem", Subsystem),
			},
			args: args{
				lsPrefix: setUpJagwPrefix("key", "", "prefix", 24),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Convert LsPrefix to Prefix no prefix",
			fields: fields{
				log: logging.DefaultLogger.WithField("subsystem", Subsystem),
			},
			args: args{
				lsPrefix: setUpJagwPrefix("key", "igpRouterId", "", 24),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Convert LsPrefix to Prefix no prefix length",
			fields: fields{
				log: logging.DefaultLogger.WithField("subsystem", Subsystem),
			},
			args: args{
				lsPrefix: setUpJagwPrefix("key", "igpRouterId", "prefix", 0),
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			adapter := &DomainAdapter{
				log: tt.fields.log,
			}
			got, err := adapter.ConvertPrefix(tt.args.lsPrefix)
			if (err != nil) != tt.wantErr {
				t.Errorf("DomainAdapter.ConvertPrefix() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if isNilInterface(got) && tt.want != nil {
				t.Errorf("DomainAdapter.ConvertPrefix() = %v, want %v", got, tt.want)
			}
			if !isNilInterface(got) && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DomainAdapter.ConvertPrefix() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDomainAdapter_ConvertPrefixEvent(t *testing.T) {
	type fields struct {
		log *logrus.Entry
	}
	tests := []struct {
		fields        fields
		name          string
		lsPrefixEvent *jagw.LsPrefixEvent
		want          domain.NetworkEvent
		wantErr       bool
	}{
		{
			name:          "nil LsPrefixEvent",
			lsPrefixEvent: nil,
			wantErr:       true,
		},
		{
			name:          "delete action success",
			lsPrefixEvent: &jagw.LsPrefixEvent{Action: proto.String("del"), Key: proto.String("key")},
			want:          domain.NewDeletePrefixEvent("key"),
			wantErr:       false,
		},
		{
			name:          "add action with ConvertPrefix error",
			lsPrefixEvent: &jagw.LsPrefixEvent{Key: proto.String("key"), Action: proto.String("add"), LsPrefix: setUpJagwPrefix("", "igpRouterId", "prefix", 24)},
			wantErr:       true,
		},
		{
			name:          "add action success",
			lsPrefixEvent: &jagw.LsPrefixEvent{Key: proto.String("key"), Action: proto.String("add"), LsPrefix: setUpJagwPrefix("key", "igpRouterId", "prefix", 24)},
			want:          domain.NewAddPrefixEvent(setUpDomainPrefix("key", "igpRouterId", "prefix", 24)),
			wantErr:       false,
		},
		{
			name:          "update action with ConvertPrefix error",
			lsPrefixEvent: &jagw.LsPrefixEvent{Key: proto.String("key"), Action: proto.String("update"), LsPrefix: setUpJagwPrefix("", "igpRouterId", "prefix", 24)},
			wantErr:       true,
		},
		{
			name:          "LsPrefixEvent with nil action",
			lsPrefixEvent: &jagw.LsPrefixEvent{Key: proto.String("key"), LsPrefix: setUpJagwPrefix("key", "igpRouterId", "prefix", 24)},
			want:          nil,
			wantErr:       true,
		},
		{
			name:          "Uknown action",
			lsPrefixEvent: &jagw.LsPrefixEvent{Key: proto.String("key"), Action: proto.String("Uknown"), LsPrefix: setUpJagwPrefix("key", "igpRouterId", "prefix", 24)},
			want:          nil,
			wantErr:       true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			adapter := &DomainAdapter{
				log: tt.fields.log,
			}
			got, err := adapter.ConvertPrefixEvent(tt.lsPrefixEvent)
			if (err != nil) != tt.wantErr {
				t.Errorf("DomainAdapter.ConvertPrefixEvent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if isNilInterface(got) && tt.want != nil {
				t.Errorf("DomainAdapter.ConvertPrefixEvent() = %v, want %v", got, tt.want)
			}
			if !isNilInterface(got) && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DomainAdapter.ConvertPrefixEvent() = %v, want %v", got, tt.want)
			}
		})
	}
}

func setUpJagwSid(key string, igpRouterId string, sid string, sidType string, algorithm uint32) *jagw.LsSrv6Sid {
	srv6Sid := &jagw.LsSrv6Sid{}
	if key != "" {
		srv6Sid.Key = proto.String(key)
	}
	if igpRouterId != "" {
		srv6Sid.IgpRouterId = proto.String(igpRouterId)
	}
	if sid != "" {
		srv6Sid.Srv6Sid = proto.String(sid)
	}
	if sidType != "" {
		srv6Sid.Srv6EndpointBehavior = &jagw.Srv6EndpointBehavior{Algorithm: proto.Uint32(algorithm)}
	}
	return srv6Sid
}

func setupDomainSid(key string, igpRouterId string, sidValue string, algorithm uint32) *domain.DomainSid {
	sid, _ := domain.NewDomainSid(&key, &igpRouterId, &sidValue, &algorithm)
	return sid
}

func TestDomainAdapter_ConvertSid(t *testing.T) {
	type fields struct {
		log *logrus.Entry
	}
	type args struct {
		lsSid *jagw.LsSrv6Sid
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    domain.Sid
		wantErr bool
	}{
		{
			name: "Convert LsSid to Sid successfully",
			fields: fields{
				log: logging.DefaultLogger.WithField("subsystem", Subsystem),
			},
			args: args{
				lsSid: setUpJagwSid("key", "igpRouterId", "sid", "sidType", 1),
			},
			want:    setupDomainSid("key", "igpRouterId", "sid", 1),
			wantErr: false,
		},
		{
			name: "Convert LsSid to Sid no key",
			fields: fields{
				log: logging.DefaultLogger.WithField("subsystem", Subsystem),
			},
			args: args{
				lsSid: setUpJagwSid("", "igpRouterId", "sid", "sidType", 1),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Convert LsSid to Sid no igp router id",
			fields: fields{
				log: logging.DefaultLogger.WithField("subsystem", Subsystem),
			},
			args: args{
				lsSid: setUpJagwSid("key", "", "sid", "sidType", 1),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Convert LsSid to Sid no sid",
			fields: fields{
				log: logging.DefaultLogger.WithField("subsystem", Subsystem),
			},
			args: args{
				lsSid: setUpJagwSid("key", "igpRouterId", "", "sidType", 1),
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			adapter := &DomainAdapter{
				log: tt.fields.log,
			}
			got, err := adapter.ConvertSid(tt.args.lsSid)
			if (err != nil) != tt.wantErr {
				t.Errorf("DomainAdapter.ConvertSid() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if isNilInterface(got) && tt.want != nil {
				t.Errorf("DomainAdapter.ConvertSid() = %v, want %v", got, tt.want)
			}
			if !isNilInterface(got) && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DomainAdapter.ConvertSid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDomainAdapter_ConvertSidEvent(t *testing.T) {
	type fields struct {
		log *logrus.Entry
	}
	tests := []struct {
		fields     fields
		name       string
		lsSidEvent *jagw.LsSrv6SidEvent
		want       domain.NetworkEvent
		wantErr    bool
	}{
		{
			name:       "nil LsSidEvent",
			lsSidEvent: nil,
			wantErr:    true,
		},
		{
			name:       "delete action success",
			lsSidEvent: &jagw.LsSrv6SidEvent{Action: proto.String("del"), Key: proto.String("key")},
			want:       domain.NewDeleteSidEvent("key"),
			wantErr:    false,
		},
		{
			name:       "add action with ConvertSid error",
			lsSidEvent: &jagw.LsSrv6SidEvent{Key: proto.String("key"), Action: proto.String("add"), LsSrv6Sid: setUpJagwSid("", "igpRouterId", "sid", "sidType", 1)},
			wantErr:    true,
		},
		{
			name:       "add action success",
			lsSidEvent: &jagw.LsSrv6SidEvent{Key: proto.String("key"), Action: proto.String("add"), LsSrv6Sid: setUpJagwSid("key", "igpRouterId", "sid", "sidType", 1)},
			want:       domain.NewAddSidEvent(setupDomainSid("key", "igpRouterId", "sid", 1)),
			wantErr:    false,
		},
		{
			name:       "update action with ConvertSid error",
			lsSidEvent: &jagw.LsSrv6SidEvent{Key: proto.String("key"), Action: proto.String("update"), LsSrv6Sid: setUpJagwSid("", "igpRouterId", "sid", "sidType", 1)},
			wantErr:    true,
		},
		{
			name:       "LsSidEvent with nil action",
			lsSidEvent: &jagw.LsSrv6SidEvent{Key: proto.String("key"), Action: nil, LsSrv6Sid: setUpJagwSid("key", "igpRouterId", "sid", "sidType", 1)},
			wantErr:    true,
		},
		{
			name:       "Uknown action",
			lsSidEvent: &jagw.LsSrv6SidEvent{Key: proto.String("key"), Action: proto.String("Uknown"), LsSrv6Sid: setUpJagwSid("key", "igpRouterId", "sid", "sidType", 1)},
			want:       nil,
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			adapter := &DomainAdapter{
				log: tt.fields.log,
			}
			got, err := adapter.ConvertSidEvent(tt.lsSidEvent)
			if (err != nil) != tt.wantErr {
				t.Errorf("DomainAdapter.ConvertSidEvent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if isNilInterface(got) && tt.want != nil {
				t.Errorf("DomainAdapter.ConvertSidEvent() = %v, want %v", got, tt.want)
			}
			if !isNilInterface(got) && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DomainAdapter.ConvertSidEvent() = %v, want %v", got, tt.want)
			}
		})
	}
}

func getDomainMinValue(value *int32) domain.Value {
	numberValue, _ := domain.NewNumberValue(domain.ValueTypeMinValue, value)
	return numberValue
}

func getDomainMaxValue(value *int32) domain.Value {
	numberValue, _ := domain.NewNumberValue(domain.ValueTypeMaxValue, value)
	return numberValue
}

func getDomainFlexAlgoValue(value *int32) domain.Value {
	numberValue, _ := domain.NewNumberValue(domain.ValueTypeFlexAlgoNr, value)
	return numberValue
}

func getDomainSfcValue(value *string) domain.Value {
	stringValue, _ := domain.NewStringValue(domain.ValueTypeSFC, value)
	return stringValue
}

func TestDomainAdapter_ConvertValuesToDomain(t *testing.T) {
	type fields struct {
		log *logrus.Entry
	}
	type args struct {
		apiValues []*api.Value
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []domain.Value
		wantErr bool
	}{
		{
			name: "Convert single min API value to domain values successfully",
			fields: fields{
				log: logging.DefaultLogger.WithField("subsystem", Subsystem),
			},
			args: args{
				apiValues: []*api.Value{
					{Type: api.ValueType_VALUE_TYPE_MIN_VALUE, NumberValue: proto.Int32(100)},
				},
			},
			want: []domain.Value{
				getDomainMinValue(proto.Int32(100)),
			},
			wantErr: false,
		},
		{
			name: "Convert single max API value to domain values successfully",
			fields: fields{
				log: logging.DefaultLogger.WithField("subsystem", Subsystem),
			},
			args: args{
				apiValues: []*api.Value{
					{Type: api.ValueType_VALUE_TYPE_MAX_VALUE, NumberValue: proto.Int32(100)},
				},
			},
			want: []domain.Value{
				getDomainMaxValue(proto.Int32(100)),
			},
			wantErr: false,
		},
		{
			name: "Convert min and max API value to domain values successfully",
			fields: fields{
				log: logging.DefaultLogger.WithField("subsystem", Subsystem),
			},
			args: args{
				apiValues: []*api.Value{
					{Type: api.ValueType_VALUE_TYPE_MAX_VALUE, NumberValue: proto.Int32(100)},
					{Type: api.ValueType_VALUE_TYPE_MIN_VALUE, NumberValue: proto.Int32(1)},
				},
			},
			want: []domain.Value{
				getDomainMaxValue(proto.Int32(100)),
				getDomainMinValue(proto.Int32(1)),
			},
			wantErr: false,
		},
		{
			name: "Convert min value to domain values error failed min < 1",
			fields: fields{
				log: logging.DefaultLogger.WithField("subsystem", Subsystem),
			},
			args: args{
				apiValues: []*api.Value{
					{Type: api.ValueType_VALUE_TYPE_MIN_VALUE, NumberValue: proto.Int32(0)},
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Convert max value to domain values error failed min < 1",
			fields: fields{
				log: logging.DefaultLogger.WithField("subsystem", Subsystem),
			},
			args: args{
				apiValues: []*api.Value{
					{Type: api.ValueType_VALUE_TYPE_MAX_VALUE, NumberValue: proto.Int32(0)},
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Convert flex algo value to domain value successfully",
			fields: fields{
				log: logging.DefaultLogger.WithField("subsystem", Subsystem),
			},
			args: args{
				apiValues: []*api.Value{
					{Type: api.ValueType_VALUE_TYPE_FLEX_ALGO_NR, NumberValue: proto.Int32(128)},
				},
			},
			want: []domain.Value{
				getDomainFlexAlgoValue(proto.Int32(128)),
			},
			wantErr: false,
		},
		{
			name: "Convert min, max and flex algo API value to domain values successfully",
			fields: fields{
				log: logging.DefaultLogger.WithField("subsystem", Subsystem),
			},
			args: args{
				apiValues: []*api.Value{
					{Type: api.ValueType_VALUE_TYPE_MAX_VALUE, NumberValue: proto.Int32(100)},
					{Type: api.ValueType_VALUE_TYPE_MIN_VALUE, NumberValue: proto.Int32(1)},
					{Type: api.ValueType_VALUE_TYPE_FLEX_ALGO_NR, NumberValue: proto.Int32(128)},
				},
			},
			want: []domain.Value{
				getDomainMaxValue(proto.Int32(100)),
				getDomainMinValue(proto.Int32(1)),
				getDomainFlexAlgoValue(proto.Int32(128)),
			},
			wantErr: false,
		},
		{
			name: "Convert string value API value to domain values successfully",
			fields: fields{
				log: logging.DefaultLogger.WithField("subsystem", Subsystem),
			},
			args: args{
				apiValues: []*api.Value{
					{Type: api.ValueType_VALUE_TYPE_SFC, StringValue: proto.String("fw")},
				},
			},
			want: []domain.Value{
				getDomainSfcValue(proto.String("fw")),
			},
			wantErr: false,
		},
		{
			name: "Convert 2 string value API value to domain values successfully",
			fields: fields{
				log: logging.DefaultLogger.WithField("subsystem", Subsystem),
			},
			args: args{
				apiValues: []*api.Value{
					{Type: api.ValueType_VALUE_TYPE_SFC, StringValue: proto.String("fw")},
					{Type: api.ValueType_VALUE_TYPE_SFC, StringValue: proto.String("ids")},
				},
			},
			want: []domain.Value{
				getDomainSfcValue(proto.String("fw")),
				getDomainSfcValue(proto.String("ids")),
			},
			wantErr: false,
		},
		{
			name: "Convert value API value to domain values - value unspecified",
			fields: fields{
				log: logging.DefaultLogger.WithField("subsystem", Subsystem),
			},
			args: args{
				apiValues: []*api.Value{
					{StringValue: proto.String("fw")},
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			adapter := &DomainAdapter{
				log: tt.fields.log,
			}
			got, err := adapter.convertValuesToDomain(tt.args.apiValues)
			if (err != nil) != tt.wantErr {
				t.Errorf("DomainAdapter.convertValuesToDomain() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if isNilInterface(got) && tt.want != nil {
				t.Errorf("DomainAdapter.convertValuesToDomain() = %v, want %v", got, tt.want)
			}
			if !isNilInterface(got) && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DomainAdapter.convertValuesToDomain() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDomainAdapter_ConvertIntentTypeToDomain(t *testing.T) {
	type fields struct {
		log *logrus.Entry
	}
	type args struct {
		apiIntentType api.IntentType
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    domain.IntentType
		wantErr bool
	}{
		{
			name: "Convert high BW API intent type to domain intent type successfully",
			fields: fields{
				log: logging.DefaultLogger.WithField("subsystem", Subsystem),
			},
			args: args{
				apiIntentType: api.IntentType_INTENT_TYPE_HIGH_BANDWIDTH,
			},
			want:    domain.IntentTypeHighBandwidth,
			wantErr: false,
		},
		{
			name: "Convert API LOW BW intent type to domain intent type successfully",
			fields: fields{
				log: logging.DefaultLogger.WithField("subsystem", Subsystem),
			},
			args: args{
				apiIntentType: api.IntentType_INTENT_TYPE_LOW_BANDWIDTH,
			},
			want:    domain.IntentTypeLowBandwidth,
			wantErr: false,
		},
		{
			name: "Convert low latency API intent type to domain intent type successfully",
			fields: fields{
				log: logging.DefaultLogger.WithField("subsystem", Subsystem),
			},
			args: args{
				apiIntentType: api.IntentType_INTENT_TYPE_LOW_LATENCY,
			},
			want:    domain.IntentTypeLowLatency,
			wantErr: false,
		},
		{
			name: "Convert low packet loss API intent type to domain intent type successfully",
			fields: fields{
				log: logging.DefaultLogger.WithField("subsystem", Subsystem),
			},
			args: args{
				apiIntentType: api.IntentType_INTENT_TYPE_LOW_PACKET_LOSS,
			},
			want:    domain.IntentTypeLowPacketLoss,
			wantErr: false,
		},
		{
			name: "Convert low jitter API intent type to domain intent type successfully",
			fields: fields{
				log: logging.DefaultLogger.WithField("subsystem", Subsystem),
			},
			args: args{
				apiIntentType: api.IntentType_INTENT_TYPE_LOW_JITTER,
			},
			want:    domain.IntentTypeLowJitter,
			wantErr: false,
		},
		{
			name: "Convert flex algo API intent type to domain intent type successfully",
			fields: fields{
				log: logging.DefaultLogger.WithField("subsystem", Subsystem),
			},
			args: args{
				apiIntentType: api.IntentType_INTENT_TYPE_FLEX_ALGO,
			},
			want:    domain.IntentTypeFlexAlgo,
			wantErr: false,
		},
		{
			name: "Convert sfc API intent type to domain intent type successfully",
			fields: fields{
				log: logging.DefaultLogger.WithField("subsystem", Subsystem),
			},
			args: args{
				apiIntentType: api.IntentType_INTENT_TYPE_SFC,
			},
			want:    domain.IntentTypeSFC,
			wantErr: false,
		},
		{
			name: "Convert low utilization API intent type to domain intent type successfully",
			fields: fields{
				log: logging.DefaultLogger.WithField("subsystem", Subsystem),
			},
			args: args{
				apiIntentType: api.IntentType_INTENT_TYPE_LOW_UTILIZATION,
			},
			want:    domain.IntentTypeLowUtilization,
			wantErr: false,
		},
		{
			name: "Convert nil API intent type to domain intent type error",
			fields: fields{
				log: logging.DefaultLogger.WithField("subsystem", Subsystem),
			},
			args:    args{},
			want:    domain.IntentTypeUnspecified,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			adapter := &DomainAdapter{
				log: tt.fields.log,
			}
			got, err := adapter.convertIntentTypeToDomain(tt.args.apiIntentType)
			if (err != nil) != tt.wantErr {
				t.Errorf("DomainAdapter.convertIntentTypeToDomain() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DomainAdapter.convertIntentTypeToDomain() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDomainAdapter_ConvertIntentsToDomain(t *testing.T) {
	type fields struct {
		log *logrus.Entry
	}
	type args struct {
		apiIntents []*api.Intent
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []domain.Intent
		wantErr bool
	}{
		{
			name: "Convert single API intent to domain intent successfully",
			fields: fields{
				log: logging.DefaultLogger.WithField("subsystem", Subsystem),
			},
			args: args{
				apiIntents: []*api.Intent{
					{
						Type: api.IntentType_INTENT_TYPE_HIGH_BANDWIDTH,
					},
				},
			},
			want: []domain.Intent{
				domain.NewDomainIntent(domain.IntentTypeHighBandwidth, []domain.Value{}),
			},
			wantErr: false,
		},
		{
			name: "Convert single API intent to domain intent value error",
			fields: fields{
				log: logging.DefaultLogger.WithField("subsystem", Subsystem),
			},
			args: args{
				apiIntents: []*api.Intent{
					{
						Type: api.IntentType_INTENT_TYPE_HIGH_BANDWIDTH,
						Values: []*api.Value{
							{Type: api.ValueType_VALUE_TYPE_MAX_VALUE, NumberValue: proto.Int32(0)},
						},
					},
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Convert single API intent to domain intent type error",
			fields: fields{
				log: logging.DefaultLogger.WithField("subsystem", Subsystem),
			},
			args: args{
				apiIntents: []*api.Intent{
					{
						Type: api.IntentType_INTENT_TYPE_UNSPECIFIED,
					},
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Convert several API intents to domain intent type error",
			fields: fields{
				log: logging.DefaultLogger.WithField("subsystem", Subsystem),
			},
			args: args{
				apiIntents: []*api.Intent{
					{
						Type: api.IntentType_INTENT_TYPE_LOW_LATENCY,
					},
					{
						Type: api.IntentType_INTENT_TYPE_LOW_PACKET_LOSS,
						Values: []*api.Value{
							{Type: api.ValueType_VALUE_TYPE_MAX_VALUE, NumberValue: proto.Int32(2)},
						},
					},
				},
			},
			want: []domain.Intent{
				domain.NewDomainIntent(domain.IntentTypeLowLatency, []domain.Value{}),
				domain.NewDomainIntent(domain.IntentTypeLowPacketLoss, []domain.Value{getDomainMaxValue(proto.Int32(2))}),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			adapter := &DomainAdapter{
				log: tt.fields.log,
			}
			got, err := adapter.convertIntentsToDomain(tt.args.apiIntents)
			if (err != nil) != tt.wantErr {
				t.Errorf("ConvertIntentsToDomain() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ConvertIntentsToDomain() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func getDomainPathRequest(source string, destination string, intents []domain.Intent, stream api.IntentController_GetIntentPathServer, ctx context.Context) domain.PathRequest {
	pathRequest, _ := domain.NewDomainPathRequest(source, destination, intents, stream, ctx)
	return pathRequest
}

func TestDomainAdapter_ConvertPathRequest(t *testing.T) {
	stream := api.NewMockIntentController_GetIntentPathServer(gomock.NewController(t))
	type fields struct {
		log *logrus.Entry
	}
	type args struct {
		pathRequest *api.PathRequest
		stream      api.IntentController_GetIntentPathServer
		ctx         context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    domain.PathRequest
		wantErr bool
	}{
		{
			name: "Convert API path request to domain path request successfully",
			fields: fields{
				log: logging.DefaultLogger.WithField("subsystem", Subsystem),
			},
			args: args{
				pathRequest: &api.PathRequest{
					Ipv6SourceAddress:      "fc:a::10",
					Ipv6DestinationAddress: "fc:b::10",
					Intents: []*api.Intent{
						{
							Type: api.IntentType_INTENT_TYPE_HIGH_BANDWIDTH,
						},
					},
				},
				stream: stream,
				ctx:    context.Background(),
			},
			want:    getDomainPathRequest("fc:a::10", "fc:b::10", []domain.Intent{domain.NewDomainIntent(domain.IntentTypeHighBandwidth, []domain.Value{})}, stream, context.Background()),
			wantErr: false,
		},
		{
			name: "Convert API path request to domain path request with values successfully",
			fields: fields{
				log: logging.DefaultLogger.WithField("subsystem", Subsystem),
			},
			args: args{
				pathRequest: &api.PathRequest{
					Ipv6SourceAddress:      "fc:a::10",
					Ipv6DestinationAddress: "fc:b::10",
					Intents: []*api.Intent{
						{
							Type: api.IntentType_INTENT_TYPE_LOW_LATENCY,
							Values: []*api.Value{
								{Type: api.ValueType_VALUE_TYPE_MAX_VALUE, NumberValue: proto.Int32(100)},
							},
						},
					},
				},
				stream: stream,
				ctx:    context.Background(),
			},
			want:    getDomainPathRequest("fc:a::10", "fc:b::10", []domain.Intent{domain.NewDomainIntent(domain.IntentTypeLowLatency, []domain.Value{getDomainMaxValue(proto.Int32(100))})}, stream, context.Background()),
			wantErr: false,
		},
		{
			name: "Convert API path request to domain path request with values error validation failed - high bandwidth does not allow max value",
			fields: fields{
				log: logging.DefaultLogger.WithField("subsystem", Subsystem),
			},
			args: args{
				pathRequest: &api.PathRequest{
					Ipv6SourceAddress:      "fc:a::10",
					Ipv6DestinationAddress: "fc:b::10",
					Intents: []*api.Intent{
						{
							Type: api.IntentType_INTENT_TYPE_HIGH_BANDWIDTH,
							Values: []*api.Value{
								{Type: api.ValueType_VALUE_TYPE_MAX_VALUE, NumberValue: proto.Int32(100)},
							},
						},
					},
				},
				stream: stream,
				ctx:    context.Background(),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Convert API path request to domain path request with values error validation failed - Undefined intent type",
			fields: fields{
				log: logging.DefaultLogger.WithField("subsystem", Subsystem),
			},
			args: args{
				pathRequest: &api.PathRequest{
					Ipv6SourceAddress:      "fc:a::10",
					Ipv6DestinationAddress: "fc:b::10",
					Intents: []*api.Intent{
						{
							Type: api.IntentType_INTENT_TYPE_UNSPECIFIED,
						},
					},
				},
				stream: stream,
				ctx:    context.Background(),
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			adapter := &DomainAdapter{
				log: tt.fields.log,
			}
			got, err := adapter.ConvertPathRequest(tt.args.pathRequest, tt.args.stream, tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("ConvertPathRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !isNilInterface(got) && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ConvertPathRequest() got = %v, want %v", got, tt.want)
			}
		})
	}
}
func TestDomainAdapter_ConvertValuesToApi(t *testing.T) {
	type fields struct {
		log *logrus.Entry
	}
	type args struct {
		values []domain.Value
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []*api.Value
	}{
		{
			name: "Convert empty domain values to API values successfully",
			fields: fields{
				log: logging.DefaultLogger.WithField("subsystem", Subsystem),
			},
			args: args{
				values: []domain.Value{},
			},
			want: []*api.Value{},
		},
		{
			name: "Convert single max domain value to API values successfully",
			fields: fields{
				log: logging.DefaultLogger.WithField("subsystem", Subsystem),
			},
			args: args{
				values: []domain.Value{
					getDomainMaxValue(proto.Int32(100)),
				},
			},
			want: []*api.Value{
				{
					Type:        api.ValueType_VALUE_TYPE_MAX_VALUE,
					NumberValue: proto.Int32(100),
				},
			},
		},
		{
			name: "Convert single min domain value to API values successfully",
			fields: fields{
				log: logging.DefaultLogger.WithField("subsystem", Subsystem),
			},
			args: args{
				values: []domain.Value{
					getDomainMinValue(proto.Int32(100)),
				},
			},
			want: []*api.Value{
				{
					Type:        api.ValueType_VALUE_TYPE_MIN_VALUE,
					NumberValue: proto.Int32(100),
				},
			},
		},
		{
			name: "Convert single flex algo domain value to API values successfully",
			fields: fields{
				log: logging.DefaultLogger.WithField("subsystem", Subsystem),
			},
			args: args{
				values: []domain.Value{
					getDomainFlexAlgoValue(proto.Int32(128)),
				},
			},
			want: []*api.Value{
				{
					Type:        api.ValueType_VALUE_TYPE_FLEX_ALGO_NR,
					NumberValue: proto.Int32(128),
				},
			},
		},
		{
			name: "Convert number domain values to API values successfully",
			fields: fields{
				log: logging.DefaultLogger.WithField("subsystem", Subsystem),
			},
			args: args{
				values: []domain.Value{
					getDomainMaxValue(proto.Int32(100)),
					getDomainMinValue(proto.Int32(10)),
				},
			},
			want: []*api.Value{
				{
					Type:        api.ValueType_VALUE_TYPE_MAX_VALUE,
					NumberValue: proto.Int32(100),
				},
				{
					Type:        api.ValueType_VALUE_TYPE_MIN_VALUE,
					NumberValue: proto.Int32(10),
				},
			},
		},
		{
			name: "Convert single sfc domain value to API values successfully",
			fields: fields{
				log: logging.DefaultLogger.WithField("subsystem", Subsystem),
			},
			args: args{
				values: []domain.Value{
					getDomainSfcValue(proto.String("fw")),
				},
			},
			want: []*api.Value{
				{
					Type:        api.ValueType_VALUE_TYPE_SFC,
					StringValue: proto.String("fw"),
				},
			},
		},
		{
			name: "Convert several sfc domain values to API values successfully",
			fields: fields{
				log: logging.DefaultLogger.WithField("subsystem", Subsystem),
			},
			args: args{
				values: []domain.Value{
					getDomainSfcValue(proto.String("fw")),
					getDomainSfcValue(proto.String("ids")),
				},
			},
			want: []*api.Value{
				{
					Type:        api.ValueType_VALUE_TYPE_SFC,
					StringValue: proto.String("fw"),
				},
				{
					Type:        api.ValueType_VALUE_TYPE_SFC,
					StringValue: proto.String("ids"),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			adapter := &DomainAdapter{
				log: tt.fields.log,
			}
			got := adapter.convertValuesToApi(tt.args.values)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ConvertValuesToApi() = %v, want %v", got, tt.want)
			}
		})
	}
}
func TestDomainAdapter_ConvertIntentsToApi(t *testing.T) {
	type fields struct {
		log *logrus.Entry
	}
	type args struct {
		intents []domain.Intent
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []*api.Intent
	}{
		{
			name: "Convert single domain intent to API intent successfully",
			fields: fields{
				log: logging.DefaultLogger.WithField("subsystem", Subsystem),
			},
			args: args{
				intents: []domain.Intent{
					domain.NewDomainIntent(domain.IntentTypeHighBandwidth, []domain.Value{}),
				},
			},
			want: []*api.Intent{
				{
					Type:   api.IntentType_INTENT_TYPE_HIGH_BANDWIDTH,
					Values: []*api.Value{},
				},
			},
		},
		{
			name: "Convert multiple domain intents to API intents successfully",
			fields: fields{
				log: logging.DefaultLogger.WithField("subsystem", Subsystem),
			},
			args: args{
				intents: []domain.Intent{
					domain.NewDomainIntent(domain.IntentTypeLowLatency, []domain.Value{}),
					domain.NewDomainIntent(domain.IntentTypeLowPacketLoss, []domain.Value{}),
				},
			},
			want: []*api.Intent{
				{
					Type:   api.IntentType_INTENT_TYPE_LOW_LATENCY,
					Values: []*api.Value{},
				},
				{
					Type:   api.IntentType_INTENT_TYPE_LOW_PACKET_LOSS,
					Values: []*api.Value{},
				},
			},
		},
		{
			name: "Convert multiple domain intents with values to API intents successfully",
			fields: fields{
				log: logging.DefaultLogger.WithField("subsystem", Subsystem),
			},
			args: args{
				intents: []domain.Intent{
					domain.NewDomainIntent(domain.IntentTypeLowLatency, []domain.Value{}),
					domain.NewDomainIntent(domain.IntentTypeLowPacketLoss, []domain.Value{getDomainMinValue(proto.Int32(10))}),
				},
			},
			want: []*api.Intent{
				{
					Type:   api.IntentType_INTENT_TYPE_LOW_LATENCY,
					Values: []*api.Value{},
				},
				{
					Type: api.IntentType_INTENT_TYPE_LOW_PACKET_LOSS,
					Values: []*api.Value{
						{
							Type:        api.ValueType_VALUE_TYPE_MIN_VALUE,
							NumberValue: proto.Int32(10),
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			adapter := &DomainAdapter{
				log: tt.fields.log,
			}
			got := adapter.convertIntentsToApi(tt.args.intents)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ConvertIntentsToApi() = %v, want %v", got, tt.want)
			}
		})
	}
}

func getDomainPathResult(ipv6SourceAddress, ipv6DestinationAddress string, ipv6SidAddresses []string, intents []domain.Intent, stream api.IntentController_GetIntentPathServer, path graph.Path) domain.PathResult {
	pathRequest := getDomainPathRequest(ipv6SourceAddress, ipv6DestinationAddress, intents, stream, context.Background())
	pathResult, _ := domain.NewDomainPathResult(pathRequest, path, ipv6SidAddresses)
	return pathResult
}

func TestDomainAdapter_ConvertPathResult(t *testing.T) {
	stream := api.NewMockIntentController_GetIntentPathServer(gomock.NewController(t))
	path := graph.NewMockPath(gomock.NewController(t))
	type fields struct {
		log *logrus.Entry
	}
	tests := []struct {
		name       string
		fields     fields
		pathResult domain.PathResult
		want       *api.PathResult
		wantErr    bool
	}{
		{
			name: "Convert domain path result to API path result successfully",
			fields: fields{
				log: logging.DefaultLogger.WithField("subsystem", Subsystem),
			},
			pathResult: getDomainPathResult("fc:a::10", "fc:b::10", []string{"fc:c::10", "fc:d::10"}, []domain.Intent{domain.NewDomainIntent(domain.IntentTypeLowLatency, []domain.Value{})}, stream, path),
			want: &api.PathResult{
				Ipv6SourceAddress:      "fc:a::10",
				Ipv6DestinationAddress: "fc:b::10",
				Ipv6SidAddresses:       []string{"fc:c::10", "fc:d::10"},
				Intents: []*api.Intent{
					{
						Type:   api.IntentType_INTENT_TYPE_LOW_LATENCY,
						Values: []*api.Value{},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Convert domain path result to API path result with values successfully",
			fields: fields{
				log: logging.DefaultLogger.WithField("subsystem", Subsystem),
			},
			pathResult: getDomainPathResult("fc:a::10", "fc:b::10", []string{"fc:c::10", "fc:d::10"}, []domain.Intent{domain.NewDomainIntent(domain.IntentTypeLowPacketLoss, []domain.Value{getDomainMaxValue(proto.Int32(10))})}, stream, path),
			want: &api.PathResult{
				Ipv6SourceAddress:      "fc:a::10",
				Ipv6DestinationAddress: "fc:b::10",
				Ipv6SidAddresses:       []string{"fc:c::10", "fc:d::10"},
				Intents: []*api.Intent{
					{
						Type: api.IntentType_INTENT_TYPE_LOW_PACKET_LOSS,
						Values: []*api.Value{
							{
								Type:        api.ValueType_VALUE_TYPE_MAX_VALUE,
								NumberValue: proto.Int32(10),
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Convert domain path result to API path result with values successfully",
			fields: fields{
				log: logging.DefaultLogger.WithField("subsystem", Subsystem),
			},
			pathResult: getDomainPathResult("fc:a::10", "fc:b::10", []string{"fc:c::10", "fc:d::10"}, []domain.Intent{domain.NewDomainIntent(domain.IntentTypeLowPacketLoss, []domain.Value{getDomainMaxValue(proto.Int32(10))}), domain.NewDomainIntent(domain.IntentTypeLowLatency, []domain.Value{})}, stream, path),
			want: &api.PathResult{
				Ipv6SourceAddress:      "fc:a::10",
				Ipv6DestinationAddress: "fc:b::10",
				Ipv6SidAddresses:       []string{"fc:c::10", "fc:d::10"},
				Intents: []*api.Intent{
					{
						Type: api.IntentType_INTENT_TYPE_LOW_PACKET_LOSS,
						Values: []*api.Value{
							{
								Type:        api.ValueType_VALUE_TYPE_MAX_VALUE,
								NumberValue: proto.Int32(10),
							},
						},
					},
					{
						Type:   api.IntentType_INTENT_TYPE_LOW_LATENCY,
						Values: []*api.Value{},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Convert domain path result - error no result found",
			fields: fields{
				log: logging.DefaultLogger.WithField("subsystem", Subsystem),
			},
			pathResult: nil,
			want:       nil,
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		adapter := &DomainAdapter{
			log: tt.fields.log,
		}
		got, err := adapter.ConvertPathResult(tt.pathResult)
		if (err != nil) != tt.wantErr {
			t.Errorf("ConvertPathResult() with name '%s' had error = %v, wantErr %v", tt.name, err, tt.wantErr)
			return
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("ConvertPathResult() '%s' = %v, want %v", tt.name, got, tt.want)
		}
	}
}
