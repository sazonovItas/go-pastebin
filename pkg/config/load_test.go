package config

import (
	"reflect"
	"testing"
	"time"
)

const testData = "testdata/"

func TestLoad(t *testing.T) {
	type testSrv struct {
		Port        string        `yaml:"port"`
		Timeout     time.Duration `yaml:"timeout"`
		IdleTimeout time.Duration `yaml:"idle_timeout" mapstructure:"idle_timeout"`
	}

	type testConfig struct {
		Name string  `yaml:"name"`
		GRPC testSrv `yaml:"grpc"`
		HTTP testSrv `yaml:"http"`
	}

	type args struct {
		cfg               any
		configPath        string
		ignoreMissingFile bool
		options           []optionFunc
	}
	tests := []struct {
		name    string
		args    args
		preconf func(t *testing.T)
		want    any
		wantErr bool
	}{
		{
			name: "Test load config from given config path",
			args: args{
				cfg:        &testConfig{},
				configPath: testData + "config.yml",
				options:    []optionFunc{WithConfigType("yaml")},
			},
			want: &testConfig{
				Name: "server",
				GRPC: testSrv{
					Port:    "8080",
					Timeout: 5 * time.Second,
				},
				HTTP: testSrv{
					Port:        "9090",
					Timeout:     5 * time.Second,
					IdleTimeout: 60 * time.Second,
				},
			},
			wantErr: false,
		},
		{
			name: "Test load config from given paths",
			args: args{
				cfg: &testConfig{},
				options: []optionFunc{
					WithConfigType("yaml"),
					WithConfigName("config"),
					WithConfigPaths([]string{testData}),
				},
			},
			want: &testConfig{
				Name: "server",
				GRPC: testSrv{
					Port:    "8080",
					Timeout: 5 * time.Second,
				},
				HTTP: testSrv{
					Port:        "9090",
					Timeout:     5 * time.Second,
					IdleTimeout: 60 * time.Second,
				},
			},
			wantErr: false,
		},
		{
			name: "Test load config from envs",
			args: args{
				cfg:               &testConfig{},
				ignoreMissingFile: true,
				options:           []optionFunc{WithEnvs("server")},
			},
			preconf: func(t *testing.T) {
				t.Setenv("SERVER_NAME", "server")
				t.Setenv("SERVER_GRPC_PORT", "8080")
				t.Setenv("SERVER_GRPC_TIMEOUT", "5s")
				t.Setenv("SERVER_HTTP_PORT", "9090")
				t.Setenv("SERVER_HTTP_TIMEOUT", "5s")
				t.Setenv("SERVER_HTTP_IDLE_TIMEOUT", "60s")
			},
			want: &testConfig{
				Name: "server",
				GRPC: testSrv{
					Port:    "8080",
					Timeout: 5 * time.Second,
				},
				HTTP: testSrv{
					Port:        "9090",
					Timeout:     5 * time.Second,
					IdleTimeout: 60 * time.Second,
				},
			},
			wantErr: false,
		},
		{
			name: "Test load config with defaults",
			args: args{
				cfg:               &testConfig{},
				ignoreMissingFile: true,
				options: []optionFunc{WithDefaults(map[string]any{
					"name":              "server",
					"grpc.port":         "8080",
					"grpc.timeout":      "5s",
					"http.port":         "9090",
					"http.timeout":      "5s",
					"http.idle_timeout": "60s",
				})},
			},
			want: &testConfig{
				Name: "server",
				GRPC: testSrv{
					Port:    "8080",
					Timeout: 5 * time.Second,
				},
				HTTP: testSrv{
					Port:        "9090",
					Timeout:     5 * time.Second,
					IdleTimeout: 60 * time.Second,
				},
			},
			wantErr: false,
		},
		{
			name: "Test load config with empty defaults",
			args: args{
				cfg:               &testConfig{},
				ignoreMissingFile: true,
			},
			want:    &testConfig{},
			wantErr: false,
		},
		{
			name: "Test load config from not exist file",
			args: args{
				cfg: &testConfig{},
				options: []optionFunc{
					WithConfigType("yaml"),
					WithConfigName("notexistence"),
					WithConfigPaths([]string{testData}),
				},
			},
			want:    &testConfig{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.preconf != nil {
				tt.preconf(t)
			}
			if err := Load(tt.args.cfg, tt.args.configPath, tt.args.ignoreMissingFile, tt.args.options...); (err != nil) != tt.wantErr {
				t.Errorf("Load() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(tt.want, tt.args.cfg) {
				t.Errorf("Load() config = %v, want = %v", tt.args.cfg, tt.want)
			}
		})
	}
}
