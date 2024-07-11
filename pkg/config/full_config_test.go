package config

import (
	"reflect"
	"testing"

	"github.com/hawkv6/hawkeye/pkg/logging"
)

func TestNewFullConfig(t *testing.T) {
	type args struct {
		jagwServiceAddress   string
		jagwRequestPort      string
		jagwSubscriptionPort string
		grpcPort             string
	}
	tests := []struct {
		name    string
		args    args
		want    *FullConfig
		wantErr bool
	}{
		{
			name: "ValidConfig",
			args: args{
				jagwServiceAddress:   "localhost",
				jagwRequestPort:      "9002",
				jagwSubscriptionPort: "9003",
				grpcPort:             "10000",
			},
			want: &FullConfig{
				BaseConfig: &BaseConfig{
					log:                logging.DefaultLogger.WithField("subsystem", Subsystem),
					jagwServiceAddress: "localhost",
					jagwRequestPort:    9002,
				},
				jagwSubscriptionPort: 9003,
				grpcPort:             10000,
			},
			wantErr: false,
		},
		{
			name: "Invalid BaseConfig - request port 0",
			args: args{
				jagwServiceAddress:   "localhost",
				jagwRequestPort:      "0",
				jagwSubscriptionPort: "9003",
				grpcPort:             "10000",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Invalid config - subscription port 0",
			args: args{
				jagwServiceAddress:   "localhost",
				jagwRequestPort:      "9002",
				jagwSubscriptionPort: "0",
				grpcPort:             "10000",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Invalid config - grpcPort port 0",
			args: args{
				jagwServiceAddress:   "localhost",
				jagwRequestPort:      "9002",
				jagwSubscriptionPort: "9003",
				grpcPort:             "0",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotConfig, err := NewFullConfig(tt.args.jagwServiceAddress, tt.args.jagwRequestPort, tt.args.jagwSubscriptionPort, tt.args.grpcPort)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewBaseConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotConfig, tt.want) {
				t.Errorf("NewBaseConfig() = %v, want %v", gotConfig, tt.want)
			}
		})
	}
}
func TestFullConfig_GetJagwServiceAddress(t *testing.T) {
	tests := []struct {
		name   string
		config *FullConfig
		want   string
	}{
		{
			name: "ValidConfig",
			config: &FullConfig{
				BaseConfig: &BaseConfig{
					jagwServiceAddress: "localhost",
				},
			},
			want: "localhost",
		},
		{
			name: "EmptyConfig",
			config: &FullConfig{
				BaseConfig: &BaseConfig{
					jagwServiceAddress: "",
				},
			},
			want: "",
		},
		{
			name: "ConfigWithIP",
			config: &FullConfig{
				BaseConfig: &BaseConfig{
					jagwServiceAddress: "192.168.0.1",
				},
			},
			want: "192.168.0.1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.config.GetJagwServiceAddress()
			if got != tt.want {
				t.Errorf("GetJagwServiceAddress() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFullConfig_GetJagwRequestPort(t *testing.T) {
	tests := []struct {
		name   string
		config *FullConfig
		want   uint16
	}{
		{
			name: "ValidConfig",
			config: &FullConfig{
				BaseConfig: &BaseConfig{
					jagwRequestPort: 9002,
				},
			},
			want: 9002,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.config.GetJagwRequestPort()
			if got != tt.want {
				t.Errorf("GetJagwRequestPort() = %v, want %v", got, tt.want)
			}
		})
	}
}
func TestFullConfig_GetJagwSubscriptionPort(t *testing.T) {
	tests := []struct {
		name   string
		config *FullConfig
		want   uint16
	}{
		{
			name: "ValidConfig",
			config: &FullConfig{
				jagwSubscriptionPort: 9003,
			},
			want: 9003,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.config.GetJagwSubscriptionPort()
			if got != tt.want {
				t.Errorf("GetJagwSubscriptionPort() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFullConfig_GetGrpcPort(t *testing.T) {
	tests := []struct {
		name   string
		config *FullConfig
		want   uint16
	}{
		{
			name: "ValidConfig",
			config: &FullConfig{
				grpcPort: 10000,
			},
			want: 10000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.config.GetGrpcPort()
			if got != tt.want {
				t.Errorf("GetGrpcPort() = %v, want %v", got, tt.want)
			}
		})
	}
}
