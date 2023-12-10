package kratox

import (
	"context"
	"encoding/json"
	"net/http"

	client "github.com/ory/kratos-client-go"

	"github.com/w6d-io/x/errorx"
	"github.com/w6d-io/x/logx"
)

// DeleteIdentity is used to delete the identity who correspond to the user id on kratos service
// if kratos is unreachable or an other issues, return nil session with statusCode of the call and error-go
func (a auth) DeleteIdentity(ctx context.Context, id string) error {
	log := logx.WithName(ctx, "DeleteIdentity")

	cfg := client.NewConfiguration()

	u, err := a.getKratosAdminAddress()
	if err != nil {
		return errorx.NewHTTP(err, http.StatusInternalServerError, "fail to get kratos address")
	}
	cfg.Scheme = u.Scheme
	cfg.Host = u.Host
	cfg.Servers = []client.ServerConfiguration{
		{
			URL: u.String(),
		},
	}

	api := client.NewAPIClient(cfg)

	r, err := api.IdentityApi.DeleteIdentity(ctx, id).Execute()
	if err != nil {
		log.Error(err, "calling fail", "name", "DeleteIdentity", "response", r)
		return errorx.NewHTTP(err, r.StatusCode, "fail to call kratos")
	}

	log.V(1).Info("identity deleted", "id", id)

	return nil
}

// UpdateIdentity is used to Update the identity with user id on kratos service
// if kratos is unreachable or an other issues, return nil session with statusCode of the call and error-go
func (a auth) UpdateIdentity(ctx context.Context, id string, schemaId string, trait map[string]interface{}) (*client.Identity, error) {
	log := logx.WithName(ctx, "UpdateIdentity")

	cfg := client.NewConfiguration()

	u, err := a.getKratosAdminAddress()
	if err != nil {
		return nil, errorx.NewHTTP(err, http.StatusInternalServerError, "fail to get kratos address")
	}
	cfg.Scheme = u.Scheme
	cfg.Host = u.Host
	cfg.Servers = []client.ServerConfiguration{
		{
			URL: u.String(),
		},
	}

	api := client.NewAPIClient(cfg)

	adminUpdateIdentityBody := *client.NewUpdateIdentityBody(
		schemaId,
		"active",
		trait,
	) // AdminUpdateIdentityBody |  (optional)

	updateIdentity, r, err := api.IdentityApi.UpdateIdentity(ctx, id).UpdateIdentityBody(adminUpdateIdentityBody).Execute()
	if err != nil {
		log.Error(err, "calling fail", "name", "UpdateIdentity", "response", r)
		return nil, errorx.NewHTTP(err, r.StatusCode, "fail to call kratos")
	}
	// response from `updateIdentity`: Identity
	log.V(1).Info("identity updated", "id", updateIdentity.Id)

	return updateIdentity, err
}

// CreateIdentity is used to create the identity with user id on kratos service
// if kratos is unreachable or an other issues, return nil session with statusCode of the call and error-go
func (a auth) CreateIdentity(ctx context.Context, schemaId string, trait map[string]interface{}) (*client.Identity, error) {
	log := logx.WithName(ctx, "CreateIdentity")

	cfg := client.NewConfiguration()

	u, err := a.getKratosAdminAddress()
	if err != nil {
		return nil, errorx.NewHTTP(err, http.StatusInternalServerError, "fail to get kratos address")
	}
	cfg.Scheme = u.Scheme
	cfg.Host = u.Host
	cfg.Servers = []client.ServerConfiguration{
		{
			URL: u.String(),
		},
	}

	api := client.NewAPIClient(cfg)

	adminCreateIdentityBody := *client.NewCreateIdentityBody(
		schemaId,
		trait,
	) // AdminCreateIdentityBody |  (optional)

	createdIdentity, r, err := api.IdentityApi.CreateIdentity(ctx).CreateIdentityBody(adminCreateIdentityBody).Execute()
	if err != nil {
		log.Error(err, "calling fail", "name", "CreateIdentity", "response", r)
		return nil, errorx.NewHTTP(err, r.StatusCode, "fail to call kratos")
	}
	// response from `AdminCreateIdentity`: Identity
	log.V(1).Info("create identity", "id", createdIdentity.Id)

	//i, err := a.Identity(ctx, r)
	return createdIdentity, err
}

// GetIdentity is used to get the identity who correspond to the user id on kratos service
// if kratos is unreachable or an other issues, return nil session with statusCode of the call and error-go
func (a auth) GetIdentity(ctx context.Context, id string) (*client.Identity, error) {
	log := logx.WithName(ctx, "GetIdentity")

	cfg := client.NewConfiguration()

	u, err := a.getKratosAdminAddress()
	if err != nil {
		return nil, errorx.NewHTTP(err, http.StatusInternalServerError, "fail to get kratos address")
	}
	cfg.Scheme = u.Scheme
	cfg.Host = u.Host
	cfg.Servers = []client.ServerConfiguration{
		{
			URL: u.String(),
		},
	}

	api := client.NewAPIClient(cfg)

	getIdentity, r, err := api.IdentityApi.GetIdentity(ctx, id).Execute()
	if err != nil {
		log.Error(err, "calling fail", "name", "GetIdentity", "response", r)
		return nil, errorx.NewHTTP(err, r.StatusCode, "fail to call kratos")
	}

	log.V(2).Info("get identity", "id", id)
	return getIdentity, err
}

// GetIdentityWithCredentials is used to get the identity who correspond to the user id on kratos service
// if kratos is unreachable or an other issues, return nil session with statusCode of the call and error-go
func (a auth) GetIdentityWithCredentials(ctx context.Context, id string) (*client.Identity, error) {
	log := logx.WithName(ctx, "GetIdentityWithCredentials")
	includeCredential := []string{"oidc"}
	cfg := client.NewConfiguration()

	u, err := a.getKratosAdminAddress()
	if err != nil {
		return nil, errorx.NewHTTP(err, http.StatusInternalServerError, "fail to get kratos address")
	}
	cfg.Scheme = u.Scheme
	cfg.Host = u.Host
	cfg.Servers = []client.ServerConfiguration{
		{
			URL: u.String(),
		},
	}

	api := client.NewAPIClient(cfg)

	getIdentity, r, err := api.IdentityApi.GetIdentity(ctx, id).IncludeCredential(includeCredential).Execute()
	if err != nil {
		log.Error(err, "calling fail", "name", "GetIdentity", "response", r)
		return nil, errorx.NewHTTP(err, r.StatusCode, "fail to call kratos")
	}

	log.V(2).Info("get identity", "id", id)
	return getIdentity, err
}

// GetIdentityFromCtx gets the session from context and retrieve the identity ID
// to make the http call
func (a auth) GetIdentityFromCtx(ctx context.Context) (*client.Identity, error) {
	sess, err := GetSessionFromCtx(ctx)
	if err != nil {
		return nil, err
	}
	return a.GetIdentity(ctx, sess.Identity.Id)
}

// GetToken returns all tokens linked with the provider
func (a auth) GetToken(ctx context.Context, providerID string) (*Provider, error) {
	log := logx.WithName(ctx, "GetTokenByHttp")
	providers, err := a.GetTokens(ctx)
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

// GetTokens returns all tokens
func (a auth) GetTokens(ctx context.Context) ([]Provider, error) {
	log := logx.WithName(ctx, "GetTokensByHttp")
	sess, err := GetSessionFromCtx(ctx)
	if err != nil {
		return nil, err
	}
	i, err := a.GetIdentityWithCredentials(ctx, sess.Identity.Id)
	if err != nil {
		return nil, err
	}

	var providers []Provider
	creds := *i.Credentials
	if cred, ok := creds[string(client.IDENTITYCREDENTIALSTYPE_OIDC)]; ok {
		if provider, ok := cred.Config["providers"]; ok {
			d, err := json.Marshal(provider)
			if err != nil {
				log.Error(err, "marshal provider failed")
				return nil, errorx.NewHTTP(err, http.StatusUnauthorized, "marshal provider failed")
			}
			if err = json.Unmarshal(d, &providers); err != nil {
				log.Error(err, "unmarshal provider failed")
				return nil, errorx.NewHTTP(err, http.StatusUnauthorized, "unmarshal provider failed")
			}
		}

	}

	return providers, nil

}

// PatchIdentity record some field of identity
func (a auth) PatchIdentity(ctx context.Context, id string, jsonPatch []client.JsonPatch) (*client.Identity, error) {
	log := logx.WithName(ctx, "PatchIdentity")
	cfg := client.NewConfiguration()

	u, err := a.getKratosAdminAddress()
	if err != nil {
		log.Error(err, "fail to get admin kratos address")
		return nil, errorx.NewHTTP(err, http.StatusInternalServerError, "fail to get kratos address")
	}
	cfg.Scheme = u.Scheme
	cfg.Host = u.Host
	cfg.Servers = []client.ServerConfiguration{
		{
			URL: u.String(),
		},
	}

	api := client.NewAPIClient(cfg)
	r := api.IdentityApi.PatchIdentity(ctx, id)
	r.JsonPatch(jsonPatch)

	i, _, err := r.Execute()
	if err != nil {
		log.Error(err, "patching failed")
		return nil, errorx.NewHTTP(err, http.StatusInternalServerError, "http patch failed")
	}
	return i, nil
}
