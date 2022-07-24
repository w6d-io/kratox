package kratos

import (
	"context"
	"net/http"
	"net/url"

	client "github.com/ory/kratos-client-go"
	"github.com/w6d-io/x/errorx"
)

type Conn struct {
	// Address is the address to the kratos micro service
	Address string `json:"address" mapstructure:"address"`
}

type Helper interface {
	// CreateIdentity is used to create the identity with user id on kratos service
	// if kratos is unreachable or an other issues, return nil session with statusCode of the call and error-go
	CreateIdentity(context.Context, string, map[string]interface{}) (*client.Identity, error)

	// GetSessionFromHTTP is used to check if the session cookie is active ( ex: session.GetActive() )
	// and also return user information
	// if session is not set, return a nil session with StatusBadRequest and error
	// if kratos is unreachable or an other issues, return nil session with statusCode of the call and error-go
	GetSessionFromHTTP(context.Context, *http.Request) (*client.Session, error)

	// GetSessionFromGRPCCtx is used to forward a session stock into a context.
	// It checks if session on context is present
	// if session is not set, return a nil session with StatusBadRequest and error
	// if kratos is unreachable or an other issues, return nil session with statusCode of the call and error-go
	GetSessionFromGRPCCtx(context.Context) (*client.Session, error)

	// GetIdentity is used to get the identity who correspond to the user id on kratos service
	// if kratos is unreachable or an other issues, return nil session with statusCode of the call and error-go
	GetIdentity(context.Context, string) (*client.Identity, error)

	// GetIdentityWithCredentials is used to get the identity who correspond to the user id on kratos service
	// if kratos is unreachable or an other issues, return nil session with statusCode of the call and error-go
	GetIdentityWithCredentials(context.Context, string) (*client.Identity, error)

	// GetIdentityFromCtx gets the session from context and retrieve the identity ID
	// to make the api call
	GetIdentityFromCtx(context.Context) (*client.Identity, error)

	// GetToken returns all tokens linked with the provider
	GetToken(context.Context, string) (*Provider, error)

	// GetTokens returns all tokens linked with the provider
	GetTokens(context.Context) ([]Provider, error)

	// UpdateIdentity is used to Update the identity with user id on kratos service
	// if kratos is unreachable or an other issues, return nil session with statusCode of the call and error-go
	// @params
	//   - context
	//   - id       : identity id
	//   - schemaId : schema id
	//   - trait
	UpdateIdentity(context.Context, string, string, map[string]interface{}) (*client.Identity, error)

	// DeleteIdentity is used to delete the identity who correspond to the user id on kratos service
	// if kratos is unreachable or an other issues, return nil session with statusCode of the call and error-go
	DeleteIdentity(context.Context, string) error
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
