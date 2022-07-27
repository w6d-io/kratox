package kratox_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	kratox "github.com/w6d-io/kratox"
	"net/http"
)

var _ = Describe("Session", func() {
	Context("GET Session", func() {
		BeforeEach(func() {
		})
		AfterEach(func() {
			kratox.Kratox = nil
		})
		It("succeeds to Get Session FromGRPCCtx", func() {
			kratox.SetAddress("127.0.0.1", 8080)
			_, err := kratox.Kratox.GetSessionFromGRPCCtx(ctx)
			Expect(err).To(nil)
		})
		It("succeeds to Get Session FromHTTP", func() {
			kratox.SetAddress("127.0.0.1", 8080)
			req := http.Request{}
			_ , err := kratox.Kratox.GetSessionFromHTTP(ctx, &req)
			Expect(err).To(nil)
		})
		})
	})