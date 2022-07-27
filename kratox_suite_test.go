package kratox_test

import (
	"context"
	kratox "github.com/w6d-io/kratox"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
	"net"
	"testing"

	zapraw "go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestKratox(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Kratox Suite")
}

var (
	ctx          context.Context
	grpcMockCli  *grpc.ClientConn
	grpcListener *bufconn.Listener
)

var _ = BeforeSuite(func() {
	encoder := zapcore.EncoderConfig{
		// Keys can be anything except the empty string.
		TimeKey:        "T",
		LevelKey:       "L",
		NameKey:        "N",
		CallerKey:      "C",
		MessageKey:     "M",
		StacktraceKey:  "S",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder,
	}
	opts := zap.Options{
		Encoder:         zapcore.NewConsoleEncoder(encoder),
		Development:     true,
		StacktraceLevel: zapcore.PanicLevel,
		Level:           zapcore.Level(int8(-2)),
	}
	ctrl.SetLogger(zap.New(zap.UseFlagOptions(&opts), zap.RawZapOpts(zapraw.AddCaller(), zapraw.AddCallerSkip(-2))))
	ctx = context.Background()
	ctx = context.WithValue(ctx, "ory_kratos_session", "test")
	grpcListener = bufconn.Listen(1024 * 1024)

	dial := func(context.Context, string) (net.Conn, error) {
		return grpcListener.Dial()
	}
	var err error
	grpcMockCli, err = grpc.DialContext(ctx, "", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithContextDialer(dial))
	Expect(err).NotTo(HaveOccurred())
}, 60)

var _ = AfterSuite(func() {
	if grpcMockCli != nil {
		_ = grpcMockCli.Close()
	}
})

type kratoxMock struct {
	kratox.Helper
	behaviour string
	token     bool
	subject   string
	code      int
	provider  string
}

func (kratoxMock) DeleteIdentity(_ context.Context, _ string) error {
	return nil
}
