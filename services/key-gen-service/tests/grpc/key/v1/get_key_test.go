//go:build e2e
// +build e2e

package key

import (
	"testing"

	keyv1 "github.com/sazonovItas/go-pastebin/services/key-gen-service/gen/pb/key/v1"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestKeyAPI(t *testing.T) {
	suite.Run(t, new(KeyClientSuite))
}

func (kcs *KeyClientSuite) TestGetKey() {
	kcs.T().Parallel()

	tests := []struct {
		name     string
		wantErr  bool
		wantCode codes.Code
	}{
		{
			name:     "Test get key api",
			wantErr:  false,
			wantCode: codes.OK,
		},
	}

	for _, tt := range tests {
		kcs.T().Run(tt.name, func(t *testing.T) {
			var err error
			if _, err = kcs.KeyClient.GetKey(&keyv1.GetKeyRequest{}); (err != nil) != tt.wantErr {
				t.Errorf("GetKey() = error %v, want %v", err, tt.wantErr)
			}

			if tt.wantCode != status.Code(err) {
				t.Errorf("GetKey() = code %d, want %d", status.Code(err), tt.wantCode)
			}
		})
	}
}
