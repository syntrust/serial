package service

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"os"
)

func loadPrivateKey(fname string) (*ecdsa.PrivateKey, error) {
	pemEncoded, err := loadPem(fname)
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

func LoadPublicKey() (*ecdsa.PublicKey, error) {
	pemEncoded, err := loadPem("keys/public.key")
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
