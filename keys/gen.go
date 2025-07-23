package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"
)

func main() {
	// Gera chave privada RSA 4096 bits
	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		panic(err)
	}

	// Salva chave privada em PEM
	privFile, err := os.Create("private.pem")
	if err != nil {
		panic(err)
	}
	defer func(privFile *os.File) {
		err := privFile.Close()
		if err != nil {
			panic(err)
		}
	}(privFile)
	err = pem.Encode(privFile, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})
	if err != nil {
		panic(err)
	}

	// Salva chave pública em PEM
	pubFile, err := os.Create("public.pem")
	if err != nil {
		panic(err)
	}
	defer func(pubFile *os.File) {
		err := pubFile.Close()
		if err != nil {
			panic(err)
		}
	}(pubFile)
	pubASN1, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		panic(err)
	}
	err = pem.Encode(pubFile, &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubASN1,
	})
	if err != nil {
		panic(err)
	}
}
