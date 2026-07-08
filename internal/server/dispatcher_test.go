package server

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"io/ioutil"
	"math/big"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestLoadKeyPair_UnencryptedKey(t *testing.T) {
	// Create temporary directory
	tmpDir, err := ioutil.TempDir("", "mmock_tls_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Generate test certificate and key
	certFile, keyFile := generateTestCertificate(t, tmpDir, "")

	// Test loading unencrypted key
	dispatcher := Dispatcher{TLSKeyPassword: ""}
	_, err = dispatcher.loadKeyPair(certFile, keyFile)
	if err != nil {
		t.Errorf("Failed to load unencrypted key pair: %v", err)
	}
}

func TestLoadKeyPair_EncryptedKey(t *testing.T) {
	// Create temporary directory
	tmpDir, err := ioutil.TempDir("", "mmock_tls_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	password := "testpassword123"

	// Generate test certificate and encrypted key
	certFile, keyFile := generateTestCertificate(t, tmpDir, password)

	// Test loading encrypted key with correct password
	dispatcher := Dispatcher{TLSKeyPassword: password}
	_, err = dispatcher.loadKeyPair(certFile, keyFile)
	if err != nil {
		t.Errorf("Failed to load encrypted key pair with correct password: %v", err)
	}

	// Test loading encrypted key with wrong password
	dispatcher.TLSKeyPassword = "wrongpassword"
	_, err = dispatcher.loadKeyPair(certFile, keyFile)
	if err == nil {
		t.Error("Expected error when loading encrypted key with wrong password")
	}

	// Test loading encrypted key with no password
	dispatcher.TLSKeyPassword = ""
	_, err = dispatcher.loadKeyPair(certFile, keyFile)
	if err == nil {
		t.Error("Expected error when loading encrypted key with no password")
	}
}

// generateTestCertificate creates a test certificate and private key for testing
func generateTestCertificate(t *testing.T, dir, password string) (certFile, keyFile string) {
	// Generate private key
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("Failed to generate private key: %v", err)
	}

	// Create certificate template
	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization:  []string{"Test Org"},
			Country:       []string{"US"},
			Province:      []string{""},
			Locality:      []string{"Test City"},
			StreetAddress: []string{""},
			PostalCode:    []string{""},
		},
		NotBefore:   time.Now(),
		NotAfter:    time.Now().Add(365 * 24 * time.Hour),
		KeyUsage:    x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		IPAddresses: nil,
	}

	// Create certificate
	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
	if err != nil {
		t.Fatalf("Failed to create certificate: %v", err)
	}

	// Write certificate file
	certFile = filepath.Join(dir, "test.crt")
	certOut, err := os.Create(certFile)
	if err != nil {
		t.Fatalf("Failed to create cert file: %v", err)
	}
	defer certOut.Close()

	if err := pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: certDER}); err != nil {
		t.Fatalf("Failed to write certificate: %v", err)
	}

	// Write private key file
	keyFile = filepath.Join(dir, "test.key")
	keyOut, err := os.Create(keyFile)
	if err != nil {
		t.Fatalf("Failed to create key file: %v", err)
	}
	defer keyOut.Close()

	// Marshal private key
	privateKeyDER, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		t.Fatalf("Failed to marshal private key: %v", err)
	}

	// Create PEM block
	block := &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: privateKeyDER,
	}

	// Encrypt if password provided
	if password != "" {
		block, err = x509.EncryptPEMBlock(rand.Reader, block.Type, block.Bytes, []byte(password), x509.PEMCipherAES256)
		if err != nil {
			t.Fatalf("Failed to encrypt private key: %v", err)
		}
	}

	if err := pem.Encode(keyOut, block); err != nil {
		t.Fatalf("Failed to write private key: %v", err)
	}

	return certFile, keyFile
}
