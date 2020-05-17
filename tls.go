package mainutil

import (
	"crypto/tls"

	"github.com/payfazz/go-errors"
)

// DefaultTLSConfig .
func DefaultTLSConfig(certfile, keyfile string) (*tls.Config, error) {
	f, err := tls.LoadX509KeyPair(certfile, keyfile)
	if err != nil {
		return nil, errors.Wrap(err)
	}
	return &tls.Config{
		Certificates: []tls.Certificate{f},

		PreferServerCipherSuites: true,
		CipherSuites: []uint16{
			// TLS 1.3
			tls.TLS_CHACHA20_POLY1305_SHA256,
			tls.TLS_AES_256_GCM_SHA384,
			tls.TLS_AES_128_GCM_SHA256,

			// TLS 1.2
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305_SHA256,
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
		},

		CurvePreferences: []tls.CurveID{
			tls.X25519,
			tls.CurveP256,
		},
	}, nil
}
