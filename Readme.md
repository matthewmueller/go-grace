# Graceful shutdown for http.Server (Go 1.8+)

## Install

```sh
go get github.com/matthewmueller/go-grace
```

## Example

This example uses httprouter, but you could use any library that implements `http.Handler`.

```go
// API struct
type API struct {
	router *httprouter.Router
}

// New API
func New() *API {
	router := httprouter.New()
	router.Handler("GET", "/", http.HandlerFunc(index))

	return &API{
		router: router,
	}
}

// Listen to a port
func (a *API) Listen(addr string) error {
	return grace.Listen(addr, a.router)
}
```

## License

MIT