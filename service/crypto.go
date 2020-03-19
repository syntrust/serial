package service

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"os"
)

func sign(infoToSign interface{}) ([]byte, []byte, error) {
	jsonValue, _ := json.Marshal(infoToSign)
	r, s, err := ecdsa.Sign(rand.Reader, privateKey, Hash(jsonValue))
	if err != nil {
		return nil, nil, err
	}
	return r.Bytes(), s.Bytes(), nil
}

func Hash(b []byte) []byte {
	h := sha256.New()
	h.Write(b)
	return h.Sum(nil)
}

func loadPrivateKey(fname string) (*ecdsa.PrivateKey, error) {
	pemEncoded, err := loadPem(fname)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(pemEncoded)
	x509Encoded := block.Bytes
	privateKey, err := x509.ParseECPrivateKey(x509Encoded)
	if err != nil {
		return nil, err
	}
	return privateKey, nil
}

func loadPem(fname string) ([]byte, error) {

	file, err := os.Open(fname)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	fileInfo, err := file.Stat()
	if err != nil {
		return nil, err
	}
	pemEncoded := make([]byte, fileInfo.Size())
	_, err = file.Read(pemEncoded)
	if err != nil {
		return nil, err
	}
	return pemEncoded, nil
}

func LoadPublicKey(fname string) (*ecdsa.PublicKey, error) {
	pemEncoded, err := loadPem(fname)
	if err != nil {
		return nil, err
	}
	blockPub, _ := pem.Decode(pemEncoded)
	x509EncodedPub := blockPub.Bytes
	genericPublicKey, err := x509.ParsePKIXPublicKey(x509EncodedPub)
	if err != nil {
		return nil, err
	}
	publicKey := genericPublicKey.(*ecdsa.PublicKey)
	return publicKey, nil
}
