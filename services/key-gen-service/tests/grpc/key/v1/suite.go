//go:build e2e
// +build e2e

package key

import (
	"os"
	"time"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/timeout"
	keyv1 "github.com/sazonovItas/go-pastebin/gen/go/pb/key/v1"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	grpcServerAddrEnv = "KEYGEN_GRPC_ADDRESS"
)

type KeyClientSuite struct {
	suite.Suite
	Address   string
	KeyClient keyv1.KeyServiceClient
	cliConn   *grpc.ClientConn
}

func (kcs *KeyClientSuite) SetupSuite() {
	if kcs.Address = os.Getenv(grpcServerAddrEnv); kcs.Address == "" {
		kcs.T().Fatal("key gen grpc server address is not specified")
	}

	cli, err := grpc.NewClient(
		kcs.Address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(timeout.UnaryClientInterceptor(5*time.Second)),
	)
	if err != nil {
		kcs.T().Fatalf("failed to connect grpc server: %v", err)
	}
	kcs.cliConn = cli

	kcs.KeyClient = keyv1.NewKeyServiceClient(cli)
}

func (kcs *KeyClientSuite) TearDownTest() {
	kcs.cliConn.Close()
}
