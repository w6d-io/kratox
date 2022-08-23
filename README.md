    # kratox

## Constants

```golang
const (
    // CookieName where is stored the cookie's session
    CookieName = "ory_kratos_session"
)
```

## Functions

### func [GetAddressFromCtx](/session.go#L127)

`func GetAddressFromCtx(ctx context.Context) (string, error)`

GetAddressFromCtx return session from context or return an error

### func [GetSessionFromCtx](/session.go#L99)

`func GetSessionFromCtx(ctx context.Context) (*client.Session, error)`

GetSessionFromCtx return session from context or return an error

### func [SetAddressDetails](/kratox.go#L132)

`func SetAddressDetails(address string, verbose bool, port ...int64)`

SetAddressDetails ip or uri and set port with verbose state. Default port is nil and default verbose is false.
In production mode is not necessary to set a verbose state in the ci configuration file

### func [SetAddressInCtx](/session.go#L141)

`func SetAddressInCtx(ctx context.Context, address string) context.Context`

SetAddressInCtx record session into context

### func [SetSessionInCtx](/session.go#L115)

`func SetSessionInCtx(ctx context.Context, session *client.Session) context.Context`

SetSessionInCtx record session into context

## Types

### type [Conn](/kratox.go#L15)

`type Conn struct { ... }`

Conn is the struct variable for connect to a kratos server

### type [ContextKey](/kratox.go#L91)

`type ContextKey int`

#### Constants

```golang
const (
    AddressKey ContextKey = iota
    SessionKey
)
```

### type [Helper](/kratox.go#L24)

`type Helper interface { ... }`

#### Variables

```golang
var (
    Kratox Helper
)
```

### type [Provider](/kratox.go#L73)

`type Provider struct { ... }`
