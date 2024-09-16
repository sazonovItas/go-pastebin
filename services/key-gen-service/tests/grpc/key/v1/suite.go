//go:build e2e
// +build e2e

package key

import (
	"os"
	"time"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/timeout"
	keyv1 "github.com/sazonovItas/go-pastebin/services/key-gen-service/gen/pb/key/v1"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	grpcServerAddrEnv = "KEYGEN_GRPC_ADDRESS"
	keyGenAPITokenEnv = "KEYGEN_CORE_API_TOKEN"
)

type KeyClientSuite struct {
	suite.Suite
	Address   string
	APIToken  string
	KeyClient keyv1.KeyServiceClient
	cliConn   *grpc.ClientConn
}

func (kcs *KeyClientSuite) SetupSuite() {
	if kcs.Address = os.Getenv(grpcServerAddrEnv); kcs.Address == "" {
		kcs.T().Fatal("key gen grpc server address is not specified")
	}

	if kcs.APIToken = os.Getenv(keyGenAPITokenEnv); kcs.APIToken == "" {
		kcs.T().Fatal("key gen api token is not specified")
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
