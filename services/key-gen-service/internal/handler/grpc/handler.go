package grpchandler

import (
	"context"

	keyv1 "github.com/sazonovItas/go-pastebin/services/key-gen-service/gen/pb/key/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type KeyGenSvc interface {
	GetKey() (string, error)
}

type keyHandler struct {
	keyGenSvc KeyGenSvc

	keyv1.UnimplementedKeyServiceServer
}

func Register(server *grpc.Server, handler keyv1.KeyServiceServer) {
	keyv1.RegisterKeyServiceServer(server, handler)
}

func NewKeyHandler(keyGenSvc KeyGenSvc) *keyHandler {
	return &keyHandler{
		keyGenSvc: keyGenSvc,
	}
}

func (kh *keyHandler) GetKey(
	ctx context.Context,
	_ *keyv1.GetKeyRequest,
) (*keyv1.GetKeyResponse, error) {
	key, err := kh.keyGenSvc.GetKey()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to generate key: %s", err.Error())
	}

	return &keyv1.GetKeyResponse{
		Msg: &keyv1.KeyMessage{
			Key: key,
		},
	}, nil
}
