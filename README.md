# kratos

## Constants

```golang
const (
    // CookieName where is stored the cookie's session
    CookieName = "ory_kratos_session"
)
```

## Types

### type [Conn](/kratos.go#L11)

`type Conn struct { ... }`

### type [Helper](/kratos.go#L16)

`type Helper interface { ... }`

### type [Kratos](/kratos.go#L39)

`type Kratos struct { ... }`

#### func (Kratos) [GetIdentityFromAPI](/identity.go#L56)

`func (k Kratos) GetIdentityFromAPI(ctx context.Context, id string) (*client.Identity, error)`

GetIdentityFromAPI is used to get the identity who correspond to the user id on kratos service
if kratos is unreachable or an other issues, return nil session with statusCode of the call and error-go

#### func (Kratos) [GetIdentityFromHTTP](/identity.go#L24)

`func (k Kratos) GetIdentityFromHTTP(ctx context.Context, id string) (*identity.Identity, error)`

GetIdentityFromHTTP is used to get the identity who correspond to the user id on kratos service
if kratos is unreachable or an other issues, return nil session with statusCode of the call and error-go

#### func (Kratos) [GetSessionFromGRPCCtx](/session.go#L46)

`func (k Kratos) GetSessionFromGRPCCtx(ctx context.Context) (*client.Session, error)`

GetSessionFromGRPCCtx is used to forward a session stock into a context.
It checks if session on context is present
if session is not set, return a nil session with StatusBadRequest and error
if kratos is unreachable or an other issues, return nil session with statusCode of the call and error-go

#### func (Kratos) [GetSessionFromHTTP](/session.go#L30)

`func (k Kratos) GetSessionFromHTTP(ctx context.Context, req *http.Request) (*client.Session, error)`

GetSessionFromHTTP is used to check if the session cookie is active ( ex: session.GetActive() )
and also return user information
if session is not set, return a nil session with StatusBadRequest and error
if kratos is unreachable or an other issues, return nil session with statusCode of the call and error-go
