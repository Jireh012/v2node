package node

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"testing"
	"time"
)

func TestPemCertificateCoversDomain(t *testing.T) {
	pemBytes := mustSelfSignedPEM(t, "l324iifio1iefni1o2qdf.561352.xyz")

	ok, err := pemCertificateCoversDomain(pemBytes, "l324iifio1iefni1o2qdf.561352.xyz")
	if err != nil || !ok {
		t.Fatalf("expected cover own domain, ok=%v err=%v", ok, err)
	}

	ok, err = pemCertificateCoversDomain(pemBytes, "l324iifio1iefni1o2qdf.ghkekrfg.821561.xyz")
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if ok {
		t.Fatal("CDN host must not be covered by cert for origin domain")
	}
}

func mustSelfSignedPEM(t *testing.T, domain string) []byte {
	t.Helper()
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	tmpl := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject:      pkix.Name{CommonName: domain},
		DNSNames:     []string{domain},
		NotBefore:    time.Now().Add(-time.Hour),
		NotAfter:     time.Now().Add(24 * time.Hour),
		KeyUsage:     x509.KeyUsageDigitalSignature,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
	}
	der, err := x509.CreateCertificate(rand.Reader, tmpl, tmpl, &key.PublicKey, key)
	if err != nil {
		t.Fatal(err)
	}
	return pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: der})
}
