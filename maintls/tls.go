package maintls

import (
	"crypto/tls"
)

// TLSConfig .
func TLSConfig() *tls.Config {
	return &tls.Config{
		PreferServerCipherSuites: true,

		MinVersion: tls.VersionTLS12,

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
			tls.CurveP521,
			tls.CurveP384,
		},
	}

}

// TLSConfigCertFile .
func TLSConfigCertFile(certfile, keyfile string) (*tls.Config, error) {
	f, err := tls.LoadX509KeyPair(certfile, keyfile)
	if err != nil {
		return nil, err
	}

	config := TLSConfig()
	config.Certificates = []tls.Certificate{f}
	return config, nil
}

// TLSConfigCertString .
func TLSConfigCertString(certpem, keypem string) (*tls.Config, error) {
	f, err := tls.X509KeyPair([]byte(certpem), []byte(keypem))
	if err != nil {
		return nil, err
	}

	config := TLSConfig()
	config.Certificates = []tls.Certificate{f}
	return config, nil
}
