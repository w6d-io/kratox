package kratos

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"

	client "github.com/ory/kratos-client-go"
	"github.com/ory/kratos/identity"
	"github.com/ory/kratos/selfservice/strategy/oidc"
	"github.com/tidwall/gjson"

	"github.com/w6d-io/x/errorx"
	"github.com/w6d-io/x/logx"
)

const (
	// path to the identities entrypoint on kratos
	identityPath = "/identities/"
)

// GetIdentityFromHTTP is used to get the identity who correspond to the user id on kratos service
// if kratos is unreachable or an other issues, return nil session with statusCode of the call and error-go
func (a auth) GetIdentityFromHTTP(ctx context.Context, id string) (*identity.Identity, error) {
	log := logx.WithName(ctx, "GetIdentityFromHTTP")

	// calling kratos at <svc>:<port>/<identityPath>/<id>
	u, err := a.getKratosAddress()
	if err != nil {
		return nil, errorx.NewHTTP(err, http.StatusInternalServerError, "unable to get kratos address")
	}
	resp, err := http.Get(u.String() + identityPath + id + "?include_credential=oidc")
	if err != nil {
		log.Error(err, "error call to kratos")
		status := http.StatusInternalServerError
		if resp != nil {
			status = resp.StatusCode
		}
		return nil, errorx.NewHTTP(err, status, "fail to call kratos")
	}

	// ready response from kratos
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(err, "error call to kratos")
		return nil, errorx.NewHTTP(err, resp.StatusCode, "fail to call kratos")
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	// unmarshalling body into identity struct
	i := identity.Identity{}
	if err := json.Unmarshal(body, &i); err != nil {
		log.Error(err, "cannot unmarshal body into identity struct")
		return nil, errorx.NewHTTP(err, http.StatusInternalServerError, "cannot unmarshal body into identity struct")
	}
	configContent := gjson.GetBytes(body, "credentials.oidc.config").String()

	if _, ok := i.Credentials[identity.CredentialsTypeOIDC]; !ok {
		i.Credentials = map[identity.CredentialsType]identity.Credentials{
			identity.CredentialsTypeOIDC: {
				Config: []byte(configContent),
			},
		}
	}
	var config oidc.CredentialsConfig
	credentials, err := i.ParseCredentials(identity.CredentialsTypeOIDC, &config)
	if err != nil {
		log.Error(err, "fail to parse credential")
		return &i, nil
	}
	i.Credentials = map[identity.CredentialsType]identity.Credentials{identity.CredentialsTypeOIDC: *credentials}
	return &i, nil
}

// GetIdentityFromCtxHTTP gets the session from context and retrieve the identity ID
// to make the http call
func (a auth) GetIdentityFromCtxHTTP(ctx context.Context) (*identity.Identity, error) {
	sess, err := GetSessionFromCtx(ctx)
	if err != nil {
		return nil, errorx.NewHTTP(err, http.StatusUnauthorized, "get session failed")
	}
	return a.GetIdentityFromHTTP(ctx, sess.Identity.Id)
}

// GetIdentityFromAPI is used to get the identity who correspond to the user id on kratos service
// if kratos is unreachable or an other issues, return nil session with statusCode of the call and error-go
func (a auth) GetIdentityFromAPI(ctx context.Context, id string) (*client.Identity, error) {
	log := logx.WithName(ctx, "GetIdentityFromAPI")

	cfg := client.NewConfiguration()
	u, err := a.getKratosAddress()
	if err != nil {
		return nil, errorx.NewHTTP(err, http.StatusInternalServerError, "fail to get kratos address")
	}
	cfg.Scheme = u.Scheme
	cfg.Host = u.Host
	api := client.NewAPIClient(cfg)

	i, rsp, err := api.V0alpha1Api.AdminGetIdentity(ctx, id).Execute()
	if err != nil {
		log.Error(err, "get identity failed")
		return nil, errorx.NewHTTP(err, rsp.StatusCode, "get identity failed")
	}
	return i, nil
}

// GetIdentityFromCtxApi gets the session from context and retrieve the identity ID
// to make the api call
func (a auth) GetIdentityFromCtxApi(ctx context.Context) (*client.Identity, error) {
	sess, err := GetSessionFromCtx(ctx)
	if err != nil {
		return nil, errorx.NewHTTP(err, http.StatusUnauthorized, "get session failed")
	}
	return a.GetIdentityFromAPI(ctx, sess.Identity.Id)
}

// GetTokenByHttp returns all tokens linked with the provider
func (a auth) GetTokenByHttp(ctx context.Context, providerID string) (*Provider, error) {
	log := logx.WithName(ctx, "GetTokenByHttp")
	providers, err := a.GetTokensByHttp(ctx)
	if err != nil {
		log.Error(err, "get all tokens failed")
		return nil, err
	}

	for _, provider := range providers {
		if provider.Provider == providerID {
			return &provider, nil
		}
	}
	logx.WithName(ctx, "GetTokenByHttp").Error(nil, "provider not match")
	return &Provider{}, nil

}

// GetTokensByHttp returns all tokens
func (a auth) GetTokensByHttp(ctx context.Context) ([]Provider, error) {
	sess, err := GetSessionFromCtx(ctx)
	log := logx.WithName(ctx, "GetTokensByHttp")

	if err != nil {
		return nil, errorx.NewHTTP(err, http.StatusUnauthorized, "get session failed")
	}

	i, err := a.GetIdentityFromHTTP(ctx, sess.Identity.Id)
	if err != nil {
		return nil, errorx.NewHTTP(err, http.StatusInternalServerError, "get identity failed")
	}

	var providers []Provider
	cred, ok := i.Credentials[identity.CredentialsTypeOIDC]
	if ok {
		gjson.GetBytes(cred.Config, "providers").ForEach(func(key, value gjson.Result) bool {
			p := &Provider{}
			err := json.Unmarshal([]byte(value.Raw), p)
			if err != nil {
				log.Error(err, "unmarshal provider failed")
				return true
			}
			providers = append(providers, *p)
			return true
		})
	}

	return providers, nil

}
