package server

import "time"

// Option is a server option.
type Option interface {
	apply(*options)
}

type options struct {
	certFile       string
	keyFile        string
	idleTimeout    time.Duration
	readTimeout    time.Duration
	writeTimeout   time.Duration
	shutdownPeriod time.Duration
}

type (
	// OptionCertFile is the server certificate file.
	OptionCertFile string

	// OptionKeyFile is the server key file.
	OptionKeyFile string

	// OptionIdleTimeout is the server idle timeout.
	OptionIdleTimeout time.Duration

	// OptionReadTimeout is the server read timeout.
	OptionReadTimeout time.Duration

	// OptionWriteTimeout is the server write timeout.
	OptionWriteTimeout time.Duration

	// OptionShutdownPeriod is the server shutdown period.
	OptionShutdownPeriod time.Duration
)

func newOptions(opts ...Option) options {
	var o options
	for _, opt := range opts {
		opt.apply(&o)
	}

	if o.idleTimeout <= 0 {
		o.idleTimeout = defaultIdleTimeout
	}

	if o.readTimeout <= 0 {
		o.readTimeout = defaultReadTimeout
	}

	if o.writeTimeout <= 0 {
		o.writeTimeout = defaultWriteTimeout
	}

	if o.shutdownPeriod <= 0 {
		o.shutdownPeriod = defaultShutdownPeriod
	}

	return o
}

func (o OptionCertFile) apply(s *options) {
	s.certFile = string(o)
}

func (o OptionKeyFile) apply(s *options) {
	s.keyFile = string(o)
}

func (o OptionIdleTimeout) apply(s *options) {
	s.idleTimeout = time.Duration(o)
}

func (o OptionReadTimeout) apply(s *options) {
	s.readTimeout = time.Duration(o)
}

func (o OptionWriteTimeout) apply(s *options) {
	s.writeTimeout = time.Duration(o)
}

func (o OptionShutdownPeriod) apply(s *options) {
	s.shutdownPeriod = time.Duration(o)
}
