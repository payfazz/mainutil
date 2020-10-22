package maintls

import (
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"encoding/hex"
	"errors"
	"strings"
)

// SetStaticPeerVerification .
func SetStaticPeerVerification(config *tls.Config, insecureSkipVerify bool, sha256SumHex ...string) {
	config.InsecureSkipVerify = insecureSkipVerify

	if len(sha256SumHex) > 0 {
		oldVerify := config.VerifyPeerCertificate

		whiteList := make(map[string]struct{})
		for _, v := range sha256SumHex {
			whiteList[strings.ToLower(v)] = struct{}{}
		}

		config.VerifyPeerCertificate = func(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {
			if oldVerify != nil {
				if err := oldVerify(rawCerts, verifiedChains); err != nil {
					return err
				}
			}

			for _, raw := range rawCerts {
				cert, err := x509.ParseCertificate(raw)
				if err != nil {
					continue
				}

				sum := sha256.Sum256(cert.RawSubjectPublicKeyInfo)
				if _, ok := whiteList[hex.EncodeToString(sum[:])]; ok {
					return nil
				}

				sum = sha256.Sum256(raw)
				if _, ok := whiteList[hex.EncodeToString(sum[:])]; ok {
					return nil
				}
			}

			return errors.New("peer cert is not in the whitelist")
		}
	}
}
