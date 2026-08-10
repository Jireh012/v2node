package node

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/wyx2685/v2node/common/file"
)

func (c *Controller) renewCertTask(_ context.Context) error {
	l, err := NewLego(c.info.Common.CertInfo)
	if err != nil {
		log.WithField("tag", c.tag).Info("new lego error: ", err)
		return nil
	}
	// SNI / CertDomain may have changed while files still exist — re-issue instead of renew.
	if ok, err := existingCertMatchesDomain(c.info.Common.CertInfo.CertFile, c.info.Common.CertInfo.CertDomain); err != nil || !ok {
		if err != nil {
			log.WithField("tag", c.tag).Info("existing cert unreadable, re-issue: ", err)
		} else {
			log.WithField("tag", c.tag).Infof("cert domain mismatch (want %s), re-issuing", c.info.Common.CertInfo.CertDomain)
		}
		_ = os.Remove(c.info.Common.CertInfo.CertFile)
		_ = os.Remove(c.info.Common.CertInfo.KeyFile)
		if err := l.CreateCert(); err != nil {
			log.WithField("tag", c.tag).Info("re-issue cert error: ", err)
		}
		return nil
	}
	err = l.RenewCert()
	if err != nil {
		log.WithField("tag", c.tag).Info("renew cert error: ", err)
		return nil
	}
	return nil
}

func (c *Controller) requestCert() error {
	cert := c.info.Common.CertInfo
	switch cert.CertMode {
	case "none", "":
	case "file":
		if cert.CertFile == "" || cert.KeyFile == "" {
			return fmt.Errorf("cert file path or key file path not exist")
		}
	case "dns", "http":
		if cert.CertFile == "" || cert.KeyFile == "" {
			return fmt.Errorf("cert file path or key file path not exist")
		}
		if file.IsExist(cert.CertFile) && file.IsExist(cert.KeyFile) {
			ok, err := existingCertMatchesDomain(cert.CertFile, cert.CertDomain)
			if err != nil {
				log.WithField("tag", c.tag).Warn("existing cert unreadable, will re-issue: ", err)
				_ = os.Remove(cert.CertFile)
				_ = os.Remove(cert.KeyFile)
			} else if ok {
				return nil
			} else {
				log.WithField("tag", c.tag).Infof(
					"cert SAN/CN does not cover server_name %q — removing old cert and re-issuing",
					cert.CertDomain)
				_ = os.Remove(cert.CertFile)
				_ = os.Remove(cert.KeyFile)
			}
		}
		l, err := NewLego(cert)
		if err != nil {
			return fmt.Errorf("create lego object error: %s", err)
		}
		err = l.CreateCert()
		if err != nil {
			return fmt.Errorf("create lego cert error: %s", err)
		}
	case "self":
		if cert.CertFile == "" || cert.KeyFile == "" {
			return fmt.Errorf("cert file path or key file path not exist")
		}
		if file.IsExist(cert.CertFile) && file.IsExist(cert.KeyFile) {
			ok, err := existingCertMatchesDomain(cert.CertFile, cert.CertDomain)
			if err == nil && ok {
				return nil
			}
			_ = os.Remove(cert.CertFile)
			_ = os.Remove(cert.KeyFile)
		}
		err := generateSelfSslCertificate(
			cert.CertDomain,
			cert.CertFile,
			cert.KeyFile)
		if err != nil {
			return fmt.Errorf("generate self cert error: %s", err)
		}
	default:
		return fmt.Errorf("unsupported certmode: %s", cert.CertMode)
	}
	return nil
}

func existingCertMatchesDomain(certFile, domain string) (bool, error) {
	if domain == "" {
		return true, nil
	}
	data, err := os.ReadFile(certFile)
	if err != nil {
		return false, err
	}
	return pemCertificateCoversDomain(data, domain)
}

// pemCertificateCoversDomain checks the leaf (first) certificate in a PEM bundle.
func pemCertificateCoversDomain(pemData []byte, domain string) (bool, error) {
	rest := pemData
	for {
		var block *pem.Block
		block, rest = pem.Decode(rest)
		if block == nil {
			break
		}
		if block.Type != "CERTIFICATE" {
			continue
		}
		cert, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			return false, err
		}
		if err := cert.VerifyHostname(domain); err != nil {
			return false, nil
		}
		return true, nil
	}
	return false, fmt.Errorf("no certificate in PEM")
}

func generateSelfSslCertificate(domain, certPath, keyPath string) error {
	key, _ := rsa.GenerateKey(rand.Reader, 2048)
	tmpl := &x509.Certificate{
		Version:      3,
		SerialNumber: big.NewInt(time.Now().Unix()),
		Subject: pkix.Name{
			CommonName: domain,
		},
		DNSNames:              []string{domain},
		BasicConstraintsValid: true,
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(30, 0, 0),
	}
	cert, err := x509.CreateCertificate(rand.Reader, tmpl, tmpl, key.Public(), key)
	if err != nil {
		return err
	}
	f, err := os.OpenFile(certPath, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	err = pem.Encode(f, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: cert,
	})
	if err != nil {
		return err
	}
	f, err = os.OpenFile(keyPath, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	err = pem.Encode(f, &pem.Block{
		Type:  "EC PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	})
	if err != nil {
		return err
	}
	return nil
}
