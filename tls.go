package mainutil

import (
	"crypto/tls"

	"github.com/payfazz/go-errors"
)

// DefaultTLSConfig .
func (env *Env) DefaultTLSConfig(certfile, keyfile string) (*tls.Config, error) {
	f, err := tls.LoadX509KeyPair(certfile, keyfile)
	if err != nil {
		return nil, errors.Wrap(err)
	}
	return &tls.Config{
		Certificates: []tls.Certificate{f},

		PreferServerCipherSuites: true,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,

			// Best disabled, as they don't provide Forward Secrecy
			// tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			// tls.TLS_RSA_WITH_AES_128_GCM_SHA256,
		},

		CurvePreferences: []tls.CurveID{
			tls.X25519,
			tls.CurveP256,
		},
	}, nil
}
