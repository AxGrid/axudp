package axcrypto

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"gitlab.com/NebulousLabs/fastrand"
)

const size = 4096

type RSAPublicKey struct {
	rsa.PublicKey
}

type RSAPrivateKey struct {
	rsa.PrivateKey
}

func CreateRSAKeyPair() (*RSAPrivateKey, *RSAPublicKey, error) {

	privateKey, err := rsa.GenerateKey(fastrand.Reader, size)
	if err != nil {
		return nil, nil, err
	}
	return &RSAPrivateKey{*privateKey}, &RSAPublicKey{privateKey.PublicKey}, nil
}

func (k *RSAPublicKey) Pem() string {
	res := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PUBLIC KEY",
			Bytes: k.Bytes(),
		},
	)
	return string(res)
}

func (k *RSAPrivateKey) Pem() string {
	res := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: k.Bytes(),
		},
	)
	return string(res)
}

func (k *RSAPublicKey) Bytes() []byte {
	return x509.MarshalPKCS1PublicKey(&k.PublicKey)
}

func (k *RSAPrivateKey) Bytes() []byte {
	return x509.MarshalPKCS1PrivateKey(&k.PrivateKey)
}

func (k *RSAPublicKey) Encrypt(message []byte) ([]byte, error) {
	encryptedBytes, err := rsa.EncryptOAEP(
		sha256.New(),
		rand.Reader,
		&k.PublicKey,
		message,
		nil)
	if err != nil {
		return nil, err
	}
	return encryptedBytes, nil
}

func (k *RSAPrivateKey) Decrypt(encryptedBytes []byte) ([]byte, error) {
	decryptedBytes, err := k.PrivateKey.Decrypt(nil, encryptedBytes, &rsa.OAEPOptions{Hash: crypto.SHA256})
	if err != nil {
		return nil, err
	}
	return decryptedBytes, nil
}

func (k *RSAPrivateKey) Sign(message []byte) ([]byte, error) {
	msgHash := sha256.New()
	_, err := msgHash.Write(message)
	if err != nil {
		return nil, err
	}
	msgHashSum := msgHash.Sum(nil)
	signature, err := rsa.SignPSS(rand.Reader, &k.PrivateKey, crypto.SHA256, msgHashSum, nil)
	if err != nil {
		return nil, err
	}
	return signature, nil
}

func (k *RSAPublicKey) Verify(message []byte, signature []byte) (bool, error) {
	msgHash := sha256.New()
	_, err := msgHash.Write(message)
	if err != nil {
		return false, err
	}
	msgHashSum := msgHash.Sum(nil)
	err = rsa.VerifyPSS(&k.PublicKey, crypto.SHA256, msgHashSum, signature, nil)
	if err != nil {
		return false, nil
	}
	return true, nil
}

func Must(b bool, err error) bool {
	if err != nil {
		return false
	}
	return b
}
