package config

import (
	"reflect"
	"testing"

	"github.com/hawkv6/hawkeye/pkg/logging"
)

func TestNewBaseConfig(t *testing.T) {
	type args struct {
		jagwServiceAddress string
		jagwRequestPort    string
	}
	tests := []struct {
		name    string
		args    args
		want    *BaseConfig
		wantErr bool
	}{
		{
			name: "Valid input localhost",
			args: args{
				jagwServiceAddress: "localhost",
				jagwRequestPort:    "8080",
			},
			want: &BaseConfig{
				log:                logging.DefaultLogger.WithField("subsystem", Subsystem),
				jagwServiceAddress: "localhost",
				jagwRequestPort:    8080,
			},
			wantErr: false,
		},
		{
			name: "Valid input 127.0.0.1",
			args: args{
				jagwServiceAddress: "127.0.0.1",
				jagwRequestPort:    "10000",
			},
			want: &BaseConfig{
				log:                logging.DefaultLogger.WithField("subsystem", Subsystem),
				jagwServiceAddress: "127.0.0.1",
				jagwRequestPort:    10000,
			},
			wantErr: false,
		},
		{
			name: "Valid input 0.0.0.0",
			args: args{
				jagwServiceAddress: "0.0.0.0",
				jagwRequestPort:    "10000",
			},
			want: &BaseConfig{
				log:                logging.DefaultLogger.WithField("subsystem", Subsystem),
				jagwServiceAddress: "0.0.0.0",
				jagwRequestPort:    10000,
			},
			wantErr: false,
		},
		{
			name: "Valid input ::",
			args: args{
				jagwServiceAddress: "::",
				jagwRequestPort:    "10000",
			},
			want: &BaseConfig{
				log:                logging.DefaultLogger.WithField("subsystem", Subsystem),
				jagwServiceAddress: "::",
				jagwRequestPort:    10000,
			},
			wantErr: false,
		},
		{
			name: "Invalid request port",
			args: args{
				jagwServiceAddress: "localhost",
				jagwRequestPort:    "0",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Invalid request port",
			args: args{
				jagwServiceAddress: "localhost",
				jagwRequestPort:    "70000",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Invalid request port",
			args: args{
				jagwServiceAddress: "localhost",
				jagwRequestPort:    "no port",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotConfig, err := NewBaseConfig(tt.args.jagwServiceAddress, tt.args.jagwRequestPort)
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

func TestGetJagwServiceAddress(t *testing.T) {
	type args struct {
		jagwServiceAddress string
		jagwRequestPort    string
	}
	tests := []struct {
		name   string
		args   args
		config *BaseConfig
		want   string
	}{
		{
			name: "Valid input",
			args: args{
				jagwServiceAddress: "localhost",
				jagwRequestPort:    "8080",
			},
			config: &BaseConfig{
				jagwServiceAddress: "localhost",
				jagwRequestPort:    8080,
			},
			want: "localhost",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := tt.config
			got := config.GetJagwServiceAddress()
			if got != tt.want {
				t.Errorf("GetJagwServiceAddress() = %v, want %v", got, tt.want)
			}
		})
	}
}
func TestGetJagwRequestPort(t *testing.T) {
	type args struct {
		jagwServiceAddress string
		jagwRequestPort    string
	}
	tests := []struct {
		name   string
		args   args
		config *BaseConfig
		want   uint16
	}{
		{
			name: "Valid input",
			args: args{
				jagwServiceAddress: "localhost",
				jagwRequestPort:    "8080",
			},
			config: &BaseConfig{
				jagwServiceAddress: "localhost",
				jagwRequestPort:    8080,
			},
			want: 8080,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := tt.config
			got := config.GetJagwRequestPort()
			if got != tt.want {
				t.Errorf("GetJagwRequestPort() = %v, want %v", got, tt.want)
			}
		})
	}
}
func TestGetJagwSubscriptionPort(t *testing.T) {
	type args struct {
		jagwServiceAddress string
		jagwRequestPort    string
	}
	tests := []struct {
		name   string
		args   args
		config *BaseConfig
		want   uint16
	}{
		{
			name: "Valid input",
			args: args{
				jagwServiceAddress: "localhost",
				jagwRequestPort:    "8080",
			},
			config: &BaseConfig{
				jagwServiceAddress: "localhost",
				jagwRequestPort:    8080,
			},
			want: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := tt.config
			got := config.GetJagwSubscriptionPort()
			if got != tt.want {
				t.Errorf("GetJagwSubscriptionPort() = %v, want %v", got, tt.want)
			}
		})
	}
}
func TestGetGrpcPort(t *testing.T) {
	type args struct {
		jagwServiceAddress string
		jagwRequestPort    string
	}
	tests := []struct {
		name   string
		args   args
		config *BaseConfig
		want   uint16
	}{
		{
			name: "Valid input",
			args: args{
				jagwServiceAddress: "localhost",
				jagwRequestPort:    "8080",
			},
			config: &BaseConfig{
				jagwServiceAddress: "localhost",
				jagwRequestPort:    8080,
			},
			want: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := tt.config
			got := config.GetGrpcPort()
			if got != tt.want {
				t.Errorf("GetGrpcPort() = %v, want %v", got, tt.want)
			}
		})
	}
}
