package kratox_test

import (
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"google.golang.org/grpc/metadata"

	kratox "github.com/w6d-io/kratox"
)

var _ = Describe("Session", func() {
	Context("GET Session", func() {
		BeforeEach(func() {
		})
		AfterEach(func() {
			kratox.Kratox = nil
		})
		It("succeeds to Get Session FromGRPCCtx", func() {
			cookieHeader := CallKratosServer("register-login")
			kratox.SetAddressDetails("127.0.0.1", Verbose, 4433)

			ctx := metadata.NewIncomingContext(ctx, metadata.MD{
				kratox.CookieName: []string{
					cookieHeader,
				},
			})
			session, err := kratox.Kratox.GetSessionFromGRPCCtx(ctx)
			Expect(err).ToNot(HaveOccurred())
			Expect(session).ToNot(BeNil())
		})
		It("error to Get Session FromGRPCCtx because cookies with name but is empty", func() {
			kratox.SetAddressDetails("127.0.0.1", Verbose, 4433)

			ctx := metadata.NewIncomingContext(ctx, metadata.MD{
				kratox.CookieName: []string{""},
			})

			session, err := kratox.Kratox.GetSessionFromGRPCCtx(ctx)
			Expect(err).To(HaveOccurred())
			Expect(session).To(BeNil())
		})
		It("error to Get Session FromGRPCCtx because cookie not found", func() {
			kratox.SetAddressDetails("127.0.0.1", Verbose, 4433)

			ctx := metadata.NewIncomingContext(ctx, metadata.MD{})
			session, err := kratox.Kratox.GetSessionFromGRPCCtx(ctx)
			Expect(err).To(HaveOccurred())
			Expect(session).To(BeNil())
		})
		It("succeeds to Get Session FromHTTP", func() {
			cookieHeader := CallKratosServer("register-login")
			kratox.SetAddressDetails("127.0.0.1", Verbose, 4433)

			req := http.Request{Header: map[string][]string{}}
			req.AddCookie(&http.Cookie{Name: kratox.CookieName, Value: cookieHeader})

			session, err := kratox.Kratox.GetSessionFromHTTP(ctx, &req)
			Expect(err).ToNot(HaveOccurred())
			Expect(session).ToNot(BeNil())
		})
		It("error to Get Session FromHTTP because cookies with name but is empty", func() {
			kratox.SetAddressDetails("127.0.0.1", Verbose, 4433)

			req := http.Request{Header: map[string][]string{}}
			req.AddCookie(&http.Cookie{Name: kratox.CookieName, Value: ""})

			session, err := kratox.Kratox.GetSessionFromHTTP(ctx, &req)
			Expect(err).To(HaveOccurred())
			Expect(session).To(BeNil())
		})
		It("error to Get Session FromHTTP because cookie not found", func() {
			kratox.SetAddressDetails("127.0.0.1", Verbose, 4433)

			req := http.Request{Header: map[string][]string{}}

			session, err := kratox.Kratox.GetSessionFromHTTP(ctx, &req)
			Expect(err).To(HaveOccurred())
			Expect(session).To(BeNil())
		})
	})
})
