package grace

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	multierror "github.com/hashicorp/go-multierror"
)

var (
	// Timeout to shutdown the server
	Timeout = 15 * time.Second
	// Signals to listen for
	Signals = []os.Signal{os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT}
)

// Listen on an addr with a handler
func Listen(addr string, handler http.Handler) error {
	errch := make(chan error)
	server := http.Server{
		Handler: handler,
		Addr:    addr,
	}

	go func() {
		errch <- server.ListenAndServe()
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, Signals...)

	var err error
	select {
	case sig := <-c:
		err = errors.New(sig.String())
	case err = <-errch:
	}

	ctx, cancel := context.WithTimeout(context.Background(), Timeout)
	defer cancel()

	if e := server.Shutdown(ctx); e != nil {
		return multierror.Append(err, e)
	}

	return err

}
