package kratos

import (
	"context"
	"encoding/json"
	"fmt"
	client "github.com/ory/kratos-client-go"
	"github.com/ory/kratos/identity"
	"github.com/ory/kratos/selfservice/strategy/oidc"
	"github.com/tidwall/gjson"
	"github.com/w6d-io/x/errorx"
	"github.com/w6d-io/x/logx"
	"io/ioutil"
	"net/http"
)

// DeleteIdentityFromHTTP is used to delete the identity who correspond to the user id on kratos service
// if kratos is unreachable or an other issues, return nil session with statusCode of the call and error-go
func (a auth) DeleteIdentityFromHTTP(ctx context.Context, id string) error {
	log := logx.WithName(ctx, "DeleteIdentityFromHTTP")

	cfg := client.NewConfiguration()

	// get kratos uriat <svc>:<port>
	u, err := a.getKratosAddress()
	if err != nil {
		return errorx.NewHTTP(err, http.StatusInternalServerError, "fail to get kratos address")
	}
	cfg.Scheme = u.Scheme
	cfg.Host = u.Host

	apiClient := client.NewAPIClient(cfg)

	r, err := apiClient.V0alpha2Api.AdminDeleteIdentity(context.Background(), id).Execute()
	if err != nil {
		log.Error(err, "Error when calling `V0alpha2Api.AdminDeleteIdentity`` %v\n")
		log.Error(err,"Full HTTP response: %v\n", r)
		return errorx.NewHTTP(err, r.StatusCode, "fail to call kratos")
	}

	logx.WithName(ctx, fmt.Sprintf("\"Successfully Removed identity\" with id %v\n", id))
	return nil
}

// UpdateIdentityFromHTTP is used to Update the identity whith user id on kratos service
// if kratos is unreachable or an other issues, return nil session with statusCode of the call and error-go
func (a auth) UpdateIdentityFromHTTP(ctx context.Context, id string, schemaId string, trait map[string]interface{}) (*identity.Identity, error) {
	log := logx.WithName(ctx, "UpdateIdentityFromHTTP")

	cfg := client.NewConfiguration()

	// get kratos uriat <svc>:<port>
	u, err := a.getKratosAddress()
	if err != nil {
		return nil, errorx.NewHTTP(err, http.StatusInternalServerError, "fail to get kratos address")
	}
	cfg.Scheme = u.Scheme
	cfg.Host = u.Host

	apiClient := client.NewAPIClient(cfg)

	adminUpdateIdentityBody  := *client.NewAdminUpdateIdentityBody(
		schemaId,
		"active",
		trait,
	) // AdminCreateIdentityBody |  (optional)

	updateIdentity, r, err := apiClient.V0alpha2Api.AdminUpdateIdentity(context.Background(), id).AdminUpdateIdentityBody(adminUpdateIdentityBody).Execute()
	if err != nil {
		log.Error(err, "Error when calling `V0alpha2Api.AdminUpdateIdentity``: %v\n")
		log.Error(err, "Full HTTP response: %v\n")
		return nil, errorx.NewHTTP(err, r.StatusCode, "fail to call kratos")
	}
	// response from `updateIdentity`: Identity
	log.Error(err, "Created identity with ID: %v\n", updateIdentity.Id)

	// ready response from kratos SDK
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error(err, "error call to kratos SDK")
		return nil, errorx.NewHTTP(err, r.StatusCode, "fail to call kratos SDK")
	}
	defer func() {
		_ = r.Body.Close()
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

// CreateIdentityFromHTTP is used to create the identity whith user id on kratos service
// if kratos is unreachable or an other issues, return nil session with statusCode of the call and error-go
func (a auth) CreateIdentityFromHTTP(ctx context.Context, schemaId string, trait map[string]interface{}) (*identity.Identity, error) {
	log := logx.WithName(ctx, "CreateIdentityFromHTTP")

	cfg := client.NewConfiguration()

	// get kratos uriat <svc>:<port>
	u, err := a.getKratosAddress()
	if err != nil {
		return nil, errorx.NewHTTP(err, http.StatusInternalServerError, "fail to get kratos address")
	}
	cfg.Scheme = u.Scheme
	cfg.Host = u.Host

	apiClient := client.NewAPIClient(cfg)

	adminCreateIdentityBody := *client.NewAdminCreateIdentityBody(
		schemaId,
		trait,
	) // AdminCreateIdentityBody |  (optional)

	createdIdentity, r, err := apiClient.V0alpha2Api.AdminCreateIdentity(context.Background()).AdminCreateIdentityBody(adminCreateIdentityBody).Execute()
	if err != nil {
		log.Error(err, "Error when calling `V0alpha2Api.AdminCreateIdentity``: %v\n")
		log.Error(err, "Full HTTP response: %v\n")
		return nil, errorx.NewHTTP(err, r.StatusCode, "fail to call kratos")
	}
	// response from `AdminCreateIdentity`: Identity
	log.Error(err, "Created identity with ID: %v\n", createdIdentity.Id)

	// ready response from kratos SDK
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error(err, "error call to kratos SDK")
		return nil, errorx.NewHTTP(err, r.StatusCode, "fail to call kratos SDK")
	}
	defer func() {
		_ = r.Body.Close()
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

// GetIdentityFromHTTP is used to get the identity who correspond to the user id on kratos service
// if kratos is unreachable or an other issues, return nil session with statusCode of the call and error-go
func (a auth) GetIdentityFromHTTP(ctx context.Context, id string) (*identity.Identity, error) {
	log := logx.WithName(ctx, "GetIdentityFromHTTP")

	cfg := client.NewConfiguration()

	// get kratos uriat <svc>:<port>
	u, err := a.getKratosAddress()
	if err != nil {
		return nil, errorx.NewHTTP(err, http.StatusInternalServerError, "fail to get kratos address")
	}
	cfg.Scheme = u.Scheme
	cfg.Host = u.Host

	apiClient := client.NewAPIClient(cfg)

	getIdentity, r, err := apiClient.V0alpha2Api.AdminGetIdentity(context.Background(), id).Execute()
	if err != nil {
		log.Error(err, "Error when calling `V0alpha2Api.AdminGetIdentity`` %v\n")
		log.Error(err,"Full HTTP response: %v\n", r)
		return nil, errorx.NewHTTP(err, r.StatusCode, "fail to call kratos")
	}

	logx.WithName(ctx, fmt.Sprintf("Data for identity with id %v. Traits %v\n", id, getIdentity.Traits))

	// ready response from kratos
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error(err, "error call to kratos")
		return nil, errorx.NewHTTP(err, r.StatusCode, "fail to call kratos")
	}
	defer func() {
		_ = r.Body.Close()
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

	i, rsp, err := api.V0alpha2Api.AdminGetIdentity(ctx, id).Execute()
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
