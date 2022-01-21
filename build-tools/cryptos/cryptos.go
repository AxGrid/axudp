package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"github.com/rs/zerolog/log"
)

func main() {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}

	priv_bytes := x509.MarshalPKCS1PrivateKey(privateKey)
	pub_bytes := x509.MarshalPKCS1PublicKey(&privateKey.PublicKey)

	pubkey_pem := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PUBLIC KEY",
			Bytes: pub_bytes,
		},
	)

	privkey_pem := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: priv_bytes,
		},
	)

	log.Info().Str("pub", string(pubkey_pem)).Str("priv", string(privkey_pem)).Msg("keys")
}
