package kratos

import (
	"context"
	"net/http"
	"net/url"

	"github.com/w6d-io/x/errorx"

	client "github.com/ory/kratos-client-go"
	"github.com/ory/kratos/identity"
)

type Conn struct {
	// Address is the address to the kratos micro service
	Address string `json:"address" mapstructure:"address"`
}

type Helper interface {

	// GetSessionFromHTTP is used to check if the session cookie is active ( ex: session.GetActive() )
	// and also return user information
	// if session is not set, return a nil session with StatusBadRequest and error
	// if kratos is unreachable or an other issues, return nil session with statusCode of the call and error-go
	GetSessionFromHTTP(ctx context.Context, req *http.Request) (*client.Session, error)

	// GetSessionFromGRPCCtx is used to forward a session stock into a context.
	// It checks if session on context is present
	// if session is not set, return a nil session with StatusBadRequest and error
	// if kratos is unreachable or an other issues, return nil session with statusCode of the call and error-go
	GetSessionFromGRPCCtx(ctx context.Context) (*client.Session, error)

	// GetIdentityFromHTTP is used to get the identity who correspond to the user id on kratos service
	// if kratos is unreachable or an other issues, return nil session with statusCode of the call and error-go
	GetIdentityFromHTTP(ctx context.Context, id string) (*identity.Identity, error)

	// GetIdentityFromAPI is used to get the identity who correspond to the user id on kratos service
	// if kratos is unreachable or an other issues, return nil session with statusCode of the call and error-go
	GetIdentityFromAPI(ctx context.Context, id string) (*client.Identity, error)

	// GetIdentityFromCtxHTTP gets the session from context and retrieve the identity ID
	// to make the http call
	GetIdentityFromCtxHTTP(ctx context.Context) (*identity.Identity, error)

	// GetIdentityFromCtxApi gets the session from context and retrieve the identity ID
	// to make the api call
	GetIdentityFromCtxApi(ctx context.Context) (*client.Identity, error)

	// GetTokenByHttp returns all tokens linked with the provider
	GetTokenByHttp(ctx context.Context, provider string) (*Provider, error)

	// GetTokensByHttp returns all tokens linked with the provider
	GetTokensByHttp(ctx context.Context) ([]Provider, error)
}

type Provider struct {
	TokenID      string `json:"initial_id_token"`
	Subject      string `json:"subject"`
	Provider     string `json:"provider"`
	AccessToken  string `json:"initial_access_token"`
	RefreshToken string `json:"initial_refresh_token"`
}

var _ Helper = &auth{}

type auth struct {
	Conn
}

var (
	Kratos Helper
)

type ContextKey int

const (
	AddressKey ContextKey = iota
	SessionKey
)

// getKratosAddress concat and format the svc and port from Conn variable
func (k Conn) getKratosAddress() (*url.URL, error) {
	u, err := url.Parse(k.Address)
	if err != nil {
		return nil, errorx.Wrap(err, "decode address failed")
	}
	if u.Host == "" {
		u, err = url.Parse("http://" + k.Address)
		if err != nil {
			return nil, errorx.Wrap(err, "decode address failed")
		}
	}
	return u, nil
}

func SetAddress(address string) {
	Kratos = &auth{Conn{Address: address}}
}
