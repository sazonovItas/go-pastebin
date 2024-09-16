//go:build e2e
// +build e2e

package key

import (
	"context"
	"testing"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/metadata"
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
		metadata metadata.MD
		wantErr  bool
		wantCode codes.Code
	}{
		{
			name:     "Test valid api token",
			metadata: metadata.MD{"x-api-key-gen-token": []string{kcs.APIToken}},
			wantErr:  false,
			wantCode: codes.OK,
		},
		{
			name:     "Test invalid api token",
			metadata: metadata.MD{"x-api-key-gen-token": []string{"invalid-api-token"}},
			wantErr:  true,
			wantCode: codes.Unauthenticated,
		},
		{
			name:     "Test invalid format api token",
			metadata: metadata.MD{"x-api-key-gen-token": []string{"bearer", "invalid-api-token"}},
			wantErr:  true,
			wantCode: codes.Unauthenticated,
		},
		{
			name:     "Test missing api token",
			metadata: metadata.MD{},
			wantErr:  true,
			wantCode: codes.Unauthenticated,
		},
	}

	for _, tt := range tests {
		kcs.T().Run(tt.name, func(t *testing.T) {
			var err error
			if _, err = kcs.KeyClient.GetKey(tt.metadata.ToOutgoing(context.Background()), &keyv1.GetKeyRequest{}); (err != nil) != tt.wantErr {
				t.Errorf("GetKey() = error %v, want %v", err, tt.wantErr)
			}

			if tt.wantCode != status.Code(err) {
				t.Errorf("GetKey() = code %d, want %d", status.Code(err), tt.wantCode)
			}
		})
	}
}
