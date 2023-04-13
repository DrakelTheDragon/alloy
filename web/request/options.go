package request

import (
	"io"
	"net/http"
)

// Option is an option for the request package.
type Option interface{ apply(*options) }

// OptionFunc is a function that implements the Option interface.
type OptionFunc func(*options)

type options struct{ byteReadLimiter byteReadLimiterFunc }

type byteReadLimiterFunc func(*http.Request) io.ReadCloser

// WithByteLimitReader sets the byteReadLimiter option which limits the number of
// bytes that can be read from the request body to the default value of 1MB.
func WithByteLimitReader(w http.ResponseWriter) OptionFunc {
	return func(o *options) {
		o.byteReadLimiter = func(r *http.Request) io.ReadCloser {
			return http.MaxBytesReader(w, r.Body, defaultMaxBytes)
		}
	}
}

func (f OptionFunc) apply(o *options) { f(o) }
