package vecsclient

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"log"
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	storeHost     = "localhost"
	storeName     = "test-store"
	storePassword = ""
	entryAlias    = "test-alias"
	entryAlias2   = "test-alias2"
)

func generatePEMAndCert() (string, string) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatal("Private key cannot be created.", err.Error())
	}

	keyPem := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	})

	r, err := rand.Int(rand.Reader, big.NewInt(1000000))
	if err != nil {
		log.Fatal("failed to generate random number")
	}
	tml := x509.Certificate{
		NotBefore:    time.Now(),
		NotAfter:     time.Now().AddDate(5, 0, 0),
		SerialNumber: r,
		Subject: pkix.Name{
			CommonName:   "VECS Go Test",
			Organization: []string{"VECS Go Test Org"},
		},
		BasicConstraintsValid: true,
	}
	cert, err := x509.CreateCertificate(rand.Reader, &tml, &tml, &key.PublicKey, key)
	if err != nil {
		log.Fatal("Certificate cannot be created.", err.Error())
	}

	certPem := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: cert,
	})

	return string(keyPem), string(certPem)
}

func TestVecsCreateStore(t *testing.T) {
	store, err := VecsCreateStore(storeHost, storeName, storePassword)
	assert.NoErrorf(t, err, "failed to create store with %+v", err)
	assert.NotNilf(t, store, "store object should not be nil")
	store.Close()

	err = VecsDeleteStore(storeHost, storeName)
	assert.NoErrorf(t, err, "failed to delete store with %+v", err)
}

func TestVecsLoadStore(t *testing.T) {
	store, err := VecsCreateStore(storeHost, storeName, storePassword)
	assert.NoErrorf(t, err, "failed to create store with %+v", err)
	assert.NotNilf(t, store, "store object should not be nil")
	store.Close()

	store, err = VecsLoadStore(storeHost, storeName, storePassword)
	assert.NoErrorf(t, err, "failed to load store with %+v", err)
	assert.NotNilf(t, store, "store object should not be nil")
	store.Close()

	err = VecsDeleteStore(storeHost, storeName)
	assert.NoErrorf(t, err, "failed to delete store with %+v", err)
}

func TestVecsDeleteStore(t *testing.T) {
	store, err := VecsCreateStore(storeHost, storeName, storePassword)
	assert.NoErrorf(t, err, "failed to create store with %+v", err)
	assert.NotNilf(t, store, "store object should not be nil")
	store.Close()

	err = VecsDeleteStore(storeHost, storeName)
	assert.NoErrorf(t, err, "failed to delete store with %+v", err)
}

func TestVecsAddEntry(t *testing.T) {
	store, err := VecsCreateStore(storeHost, storeName, storePassword)
	assert.NoErrorf(t, err, "failed to create store with %+v", err)
	assert.NotNilf(t, store, "store object should not be nil")

	pem, cert := generatePEMAndCert()

	err = store.AddEntry(VecsTypePrivateKey, entryAlias, cert, pem, storePassword, false)
	assert.NoErrorf(t, err, "failed to add entry to store with %+v", err)

	store.Close()
	err = VecsDeleteStore(storeHost, storeName)
	assert.NoErrorf(t, err, "failed to delete store with %+v", err)
}

func TestVecsGetEntry(t *testing.T) {
	store, err := VecsCreateStore(storeHost, storeName, storePassword)
	assert.NoErrorf(t, err, "failed to create store with %+v", err)
	assert.NotNilf(t, store, "store object should not be nil")

	pem, cert := generatePEMAndCert()

	err = store.AddEntry(VecsTypePrivateKey, entryAlias, cert, pem, storePassword, false)
	assert.NoErrorf(t, err, "failed to add entry to store with %+v", err)

	entry, err := store.GetEntry(entryAlias, VecsEntryInfoLevel2)
	assert.NoErrorf(t, err, "failed to get entry from store with %+v", err)
	assert.NotNilf(t, entry, "entry object should not be nil")

	entry.Free()
	store.Close()
	err = VecsDeleteStore(storeHost, storeName)
	assert.NoErrorf(t, err, "failed to delete store with %+v", err)
}

func TestVecsDeleteEntry(t *testing.T) {
	store, err := VecsCreateStore(storeHost, storeName, storePassword)
	assert.NoErrorf(t, err, "failed to create store with %+v", err)
	assert.NotNilf(t, store, "store object should not be nil")

	pem, cert := generatePEMAndCert()

	err = store.AddEntry(VecsTypePrivateKey, entryAlias, cert, pem, storePassword, false)
	assert.NoErrorf(t, err, "failed to add entry to store with %+v", err)

	err = store.DeleteEntry(entryAlias)
	assert.NoErrorf(t, err, "failed to get entry from store with %+v", err)

	store.Close()
	err = VecsDeleteStore(storeHost, storeName)
	assert.NoErrorf(t, err, "failed to delete store with %+v", err)
}

func TestVecsGetEntryCount(t *testing.T) {
	store, err := VecsCreateStore(storeHost, storeName, storePassword)
	assert.NoErrorf(t, err, "failed to create store with %+v", err)
	assert.NotNilf(t, store, "store object should not be nil")

	count, err := store.EntryCount()
	assert.NoErrorf(t, err, "failed to get entry count with %+v", err)
	assert.Equalf(t, 0, count, "entry count should be 0 but was %d", count)

	pem, cert := generatePEMAndCert()
	pem2, cert2 := generatePEMAndCert()

	err = store.AddEntry(VecsTypePrivateKey, entryAlias, cert, pem, storePassword, false)
	assert.NoErrorf(t, err, "failed to add entry to store with %+v", err)
	count, err = store.EntryCount()
	assert.NoErrorf(t, err, "failed to get entry count with %+v", err)
	assert.Equalf(t, 1, count, "entry count should be 1 but was %d", count)

	err = store.AddEntry(VecsTypePrivateKey, entryAlias2, cert2, pem2, storePassword, false)
	assert.NoErrorf(t, err, "failed to add entry to store with %+v", err)
	count, err = store.EntryCount()
	assert.NoErrorf(t, err, "failed to get entry count with %+v", err)
	assert.Equalf(t, 2, count, "entry count should be 2 but was %d", count)

	err = store.DeleteEntry(entryAlias)
	assert.NoErrorf(t, err, "failed to get entry from store with %+v", err)
	count, err = store.EntryCount()
	assert.NoErrorf(t, err, "failed to get entry count with %+v", err)
	assert.Equalf(t, 1, count, "entry count should be 1 but was %d", count)

	err = store.DeleteEntry(entryAlias2)
	assert.NoErrorf(t, err, "failed to get entry from store with %+v", err)
	count, err = store.EntryCount()
	assert.NoErrorf(t, err, "failed to get entry count with %+v", err)
	assert.Equalf(t, 0, count, "entry count should be 0 but was %d", count)

	store.Close()
	err = VecsDeleteStore(storeHost, storeName)
	assert.NoErrorf(t, err, "failed to delete store with %+v", err)
}

func TestVecsGetEntryCert(t *testing.T) {
	store, err := VecsCreateStore(storeHost, storeName, storePassword)
	assert.NoErrorf(t, err, "failed to create store with %+v", err)
	assert.NotNilf(t, store, "store object should not be nil")

	pem, cert := generatePEMAndCert()

	err = store.AddEntry(VecsTypePrivateKey, entryAlias, cert, pem, storePassword, false)
	assert.NoErrorf(t, err, "failed to add entry to store with %+v", err)

	entryCert, err := store.GetEntryCert(entryAlias)
	assert.NoErrorf(t, err, "failed to get entry cert from store with %+v", err)
	assert.Equalf(t, cert, entryCert, "generated cert does not equal added cert")

	store.Close()
	err = VecsDeleteStore(storeHost, storeName)
	assert.NoErrorf(t, err, "failed to delete store with %+v", err)
}

func TestVecsGetEntryPEM(t *testing.T) {
	store, err := VecsCreateStore(storeHost, storeName, storePassword)
	assert.NoErrorf(t, err, "failed to create store with %+v", err)
	assert.NotNilf(t, store, "store object should not be nil")

	pemStr, cert := generatePEMAndCert()

	err = store.AddEntry(VecsTypePrivateKey, entryAlias, cert, pemStr, storePassword, false)
	assert.NoErrorf(t, err, "failed to add entry to store with %+v", err)

	entryPEM, err := store.GetEntryKey(entryAlias, storePassword)
	assert.NoErrorf(t, err, "failed to get entry cert from store with %+v", err)

	// Convert to PKCS1 format and encode
	block, _ := pem.Decode([]byte(entryPEM))
	parseResult, _ := x509.ParsePKCS8PrivateKey(block.Bytes)
	key := parseResult.(*rsa.PrivateKey)
	keyPem := string(pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	}))

	assert.Equalf(t, pemStr, keyPem, "generated PEM does not equal added PEM")

	store.Close()
	err = VecsDeleteStore(storeHost, storeName)
	assert.NoErrorf(t, err, "failed to delete store with %+v", err)
}