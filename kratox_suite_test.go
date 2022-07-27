package kratox_test

import (
	"bytes"
	"context"
	"fmt"
	"github.com/jaswdr/faker"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	_ "github.com/ory/kratos-client-go"
	"github.com/w6d-io/kratox"
	zapraw "go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
	"net"
	"os/exec"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"testing"
)

func TestKratox(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Kratox Suite")
}

func MockKratosServer(arg string) {
	e := exec.Command("make", arg)
	var out bytes.Buffer
	e.Stdout = &out
	err := e.Run()
	if err != nil {
		fmt.Printf("Error %s: %q\n",arg, err)
	}
}

const (
	// Verbose state to call Kratos serveer
	Verbose = false
)

func CallKratosServer(flow string) string {
	// Ramdom faker
	faker := faker.New()
	name := faker.Person().FirstName()
	password := faker.RandomStringWithLength(16)

	// run .sh script with option and return the set-cookies header
	e := exec.Command("./makefile.sh", flow, name, password)
	var out bytes.Buffer
	e.Stdout = &out

	err := e.Run()
	if err != nil {
		fmt.Printf("Error %s: %q\n",flow, err)
	}

	// display the cookie
	fmt.Printf(kratox.CookieName + " : %q\n", out.String())

	return out.String()
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

	// start mock kratos
	MockKratosServer("kratos")
}, 60)

var _ = AfterSuite(func() {
	// stop mock kratos
	// Comment out this feature to re-run the test faster in a local environment.
	// Because the execution of the command to stop the KRATOS server process can be very long.
	MockKratosServer("stop")

	if grpcMockCli != nil {
		_ = grpcMockCli.Close()
	}
})
