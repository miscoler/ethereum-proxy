package graceful

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/pkg/errors"
)

func ListenAndServe(addr string, handler http.Handler) error {
	srv := http.Server{ //nolint:exhaustivestruct
		Addr:    addr,
		Handler: handler,
	}

	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		if err := srv.Shutdown(context.Background()); err != nil {
			log.Printf("HTTP server Shutdown: %v", err)
		}
		close(idleConnsClosed)
	}()

	err := srv.ListenAndServe()
	<-idleConnsClosed

	return errors.Wrap(err, "gracefu;ly serving http")
}
