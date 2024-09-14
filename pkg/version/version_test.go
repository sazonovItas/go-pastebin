package versioninfo

import (
	"reflect"
	"testing"

	goversion "github.com/caarlos0/go-version"
)

func TestBuildVersion(t *testing.T) {
	type args struct {
		name    string
		desc    string
		version string
		commit  string
		date    string
	}
	tests := []struct {
		name string
		args args
		want goversion.Info
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BuildVersion(tt.args.name, tt.args.desc, tt.args.version, tt.args.commit, tt.args.date); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BuildVersion() = %v, want %v", got, tt.want)
			}
		})
	}
}
