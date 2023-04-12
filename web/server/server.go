// Package server provides a simple web server.
package server

import (
	"context"
	"crypto/tls"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	defaultIdleTimeout    = time.Minute
	defaultReadTimeout    = 10 * time.Second
	defaultWriteTimeout   = 30 * time.Second
	defaultShutdownPeriod = 20 * time.Second
)

// Run starts a web server.
func Run(addr string, h http.Handler, opts ...Option) error {
	o := newOptions(opts...)

	srv := &http.Server{
		Addr:         addr,
		Handler:      h,
		IdleTimeout:  o.idleTimeout,
		ReadTimeout:  o.readTimeout,
		WriteTimeout: o.writeTimeout,
	}

	shutdownError := make(chan error)

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit

		ctx, cancel := context.WithTimeout(context.Background(), o.shutdownPeriod)
		defer cancel()

		shutdownError <- srv.Shutdown(ctx)
	}()

	var err error

	if o.certFile != "" && o.keyFile != "" {
		err = runTLS(srv, o.certFile, o.keyFile)
	} else {
		err = run(srv)
	}

	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return <-shutdownError
}

func run(srv *http.Server) error {
	return srv.ListenAndServe()
}

func runTLS(srv *http.Server, certFile, keyFile string) error {
	tlsConfig := &tls.Config{
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
		MinVersion:       tls.VersionTLS13,
	}

	srv.TLSConfig = tlsConfig

	return srv.ListenAndServeTLS(certFile, keyFile)
}
