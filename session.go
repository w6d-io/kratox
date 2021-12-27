package kratos

import (
	"context"
	"fmt"
	"net/http"

	client "github.com/ory/kratos-client-go"
	"github.com/pkg/errors"
	"google.golang.org/grpc/metadata"

	"github.com/w6d-io/x/errorx"
	"github.com/w6d-io/x/logx"
)

const (
	// CookieName where is stored the cookie's session
	CookieName = "ory_kratos_session"
)

var (
	errNoCookie             = errors.New(CookieName + " cookie not found")
	errNoMDFromCtx          = errors.New("cannot get metadata from context")
	errSessNotFoundInCtx    = errors.New("session not found in context")
	errAddressNotFoundInCtx = errors.New("address not found in context")
)

// GetSessionFromHTTP is used to check if the session cookie is active ( ex: session.GetActive() )
// and also return user information
// if session is not set, return a nil session with StatusBadRequest and error
// if kratos is unreachable or an other issues, return nil session with statusCode of the call and error-go
func (a auth) GetSessionFromHTTP(ctx context.Context, req *http.Request) (*client.Session, error) {
	log := logx.WithName(ctx, "GetSessionFromHTTP")

	cookie, err := req.Cookie(CookieName)
	if err != nil {
		log.Error(err, "get ory_kratos_session from cookie failed")
		return nil, errorx.NewHTTP(err, http.StatusBadRequest, "get ory_kratos_session from cookie failed")
	}
	log.V(2).Info("cookie", "value", cookie.Value)
	return a.do(ctx, cookie.Value)
}

// GetSessionFromGRPCCtx is used to forward a session stock into a context.
// It checks if session on context is present
// if session is not set, return a nil session with StatusBadRequest and error
// if kratos is unreachable or an other issues, return nil session with statusCode of the call and error-go
func (a auth) GetSessionFromGRPCCtx(ctx context.Context) (*client.Session, error) {
	log := logx.WithName(ctx, "GetSessionFromGRPCCtx")

	//get metadata from ctx
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		log.Error(errNoMDFromCtx, "metadata boolean from metadata.FromIncomingContext(ctx) = %v", ok)
		return nil, errorx.NewHTTP(errNoMDFromCtx, http.StatusNotFound, "fail to get metadata")
	}

	// check if session is present on our metadata
	if _, ok := md[CookieName]; !ok {
		log.Error(errNoCookie, `metadata "%s" doesn't exist`, CookieName)
		return nil, errorx.NewHTTP(errNoCookie, http.StatusNotFound, "bad metadata")
	}

	// check if we have more than zero value for this key cause MD is map[string][]string
	if len(md[CookieName]) == 0 || len(md[CookieName][0]) == 0 {
		log.Error(errNoCookie, "metadata \"%s\" exist but no value exist", CookieName)
		return nil, errorx.NewHTTP(errNoCookie, http.StatusNotFound, "empty metadata")
	}
	return a.do(ctx, md[CookieName][0])
}

func (a auth) do(ctx context.Context, cookie string) (*client.Session, error) {
	log := logx.WithName(ctx, "GetSessionFromCtx")
	u, err := a.getKratosAddress()
	if err != nil {
		return nil, errorx.NewHTTP(err, http.StatusInternalServerError, "fail to get kratos address")
	}
	cfg := client.NewConfiguration()
	cfg.Scheme = u.Scheme
	cfg.Host = u.Host
	api := client.NewAPIClient(cfg)
	log.V(2).Info("making call to kratos.GetSession", "session_id", cookie)

	sess, rsp, err := api.V0alpha1Api.ToSession(ctx).Cookie(fmt.Sprintf("%s=%s", CookieName, cookie)).Execute()
	if err != nil {
		log.Error(err, "get session failed")
		status := http.StatusInternalServerError
		if rsp != nil {
			status = rsp.StatusCode
		}
		return nil, errorx.NewHTTP(err, status, "get session failed")
	}
	return sess, nil
}

// GetSessionFromCtx return session from context or return an error
func GetSessionFromCtx(ctx context.Context) (*client.Session, error) {
	s := ctx.Value(SessionKey)
	if s == nil {
		return nil, errSessNotFoundInCtx
	}
	sess, ok := s.(*client.Session)
	if !ok {
		return nil, errSessNotFoundInCtx
	}
	if sess == nil {
		return nil, errSessNotFoundInCtx
	}
	return sess, nil

}

// SetSessionInCtx record session into context
func SetSessionInCtx(ctx context.Context, session *client.Session) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	if session == nil {
		return ctx
	}
	ctx = context.WithValue(ctx, SessionKey, session)
	return ctx
}

// GetAddressFromCtx return session from context or return an error
func GetAddressFromCtx(ctx context.Context) (string, error) {
	a := ctx.Value(AddressKey)
	if a == nil {
		return "", errAddressNotFoundInCtx
	}
	address, ok := a.(string)
	if !ok {
		return "", errAddressNotFoundInCtx
	}
	return address, nil

}

// SetAddressInCtx record session into context
func SetAddressInCtx(ctx context.Context, address string) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	if address == "" {
		return ctx
	}
	ctx = context.WithValue(ctx, AddressKey, address)
	return ctx
}
