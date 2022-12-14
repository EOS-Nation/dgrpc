package dgrpc_test

import (
	"fmt"
	"os"
	"time"

	"github.com/streamingfast/dgrpc"
	"github.com/streamingfast/logging"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func init() {
	logging.ApplicationLogger("example", "github.com/streamingfast/dgrpc_example_secure_server", &zlog)
}

func ExampleNewServer_Secure() {
	secureConfig, err := dgrpc.SecuredByX509KeyPair("./example/cert/cert.pem", "./example/cert/key.pem")
	if err != nil {
		panic(fmt.Errorf("unable to create X509 secure config: %w", err))
	}

	server := dgrpc.NewServer2(
		dgrpc.SecureServer(secureConfig),
		dgrpc.WithLogger(zlog),
		dgrpc.WithHealthCheck(dgrpc.HealthCheckOverGRPC, healthCheck),
	)

	server.RegisterService(func(gs *grpc.Server) {
		// Register some more gRPC services here against `gs`
		// pbstatedb.RegisterStateService(gs, implementation)
	})

	server.OnTerminated(func(err error) {
		if err != nil {
			zlog.Error("gRPC server unexpected failure", zap.Error(err))
		}

		// Should be tied to application lifecycle to avoid abrupt tear down
		zlog.Core().Sync()
		os.Exit(1)
	})

	go server.Launch("localhost:9000")

	// We wait 5m before shutting down, in reality you would tie that so lifecycle of your app
	time.Sleep(5 * time.Minute)

	// Gives 30s for a gracefull shutdown
	server.Shutdown(30 * time.Second)
}
