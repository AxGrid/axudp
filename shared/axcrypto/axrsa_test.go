package axcrypto

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateRSAKeyPair(t *testing.T) {
	pub, priv, err := CreateRSAKeyPair()
	assert.Nil(t, err)
	assert.NotNil(t, pub)
	assert.NotNil(t, priv)
}

func TestRSA_Pem(t *testing.T) {
	priv, pub, err := CreateRSAKeyPair()
	assert.Nil(t, err)
	pub_pem := pub.Pem()
	priv_pem := priv.Pem()
	assert.NotNil(t, pub_pem)
	assert.NotNil(t, priv_pem)
	assert.Equal(t, "-----BEGIN RSA PRIVATE KEY-----", priv_pem[:len("-----BEGIN RSA PRIVATE KEY-----")])
	assert.Equal(t, "-----BEGIN RSA PUBLIC KEY-----", pub_pem[:len("-----BEGIN RSA PUBLIC KEY-----")])
}

func TestRSA_Sign(t *testing.T) {
	priv, pub, err := CreateRSAKeyPair()
	assert.Nil(t, err)
	msg := []byte("hello from rsa")
	sign, err := priv.Sign(msg)
	assert.Nil(t, err)
	ver, err := pub.Verify(msg, sign)
	assert.Nil(t, err)
	assert.True(t, ver)
	ver, err = pub.Verify([]byte("hello from dsa"), sign)
	assert.Nil(t, err)
	assert.False(t, ver)
}

func TestRSA_Crypt(t *testing.T) {
	priv, pub, err := CreateRSAKeyPair()
	assert.Nil(t, err)
	msg := []byte("hello from rsa")
	crypted, err := pub.Encrypt(msg)
	assert.Nil(t, err)
	assert.NotEqual(t, crypted, msg)
	encrypted, err := priv.Decrypt(crypted)
	assert.Nil(t, err)
	assert.Equal(t, encrypted, msg)
}
