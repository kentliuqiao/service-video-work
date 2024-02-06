package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
	"os"
)

func main() {

	err := genkey()
	if err != nil {
		log.Fatalln(err)
	}

}

func genkey() error {

	// Generate a new key pair
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return fmt.Errorf("generating private key: %w", err)
	}

	// create a file for the private key in pem format
	privateKeyFile, err := os.Create("private.pem")
	if err != nil {
		return fmt.Errorf("creating private key file: %w", err)
	}
	defer privateKeyFile.Close()

	// construct a pem block for the private key
	privateKeyBlock := &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}

	// write the pem block to the file
	err = pem.Encode(privateKeyFile, privateKeyBlock)
	if err != nil {
		return fmt.Errorf("encoding private key to file: %w", err)
	}

	// create a file for the public key in pem format
	publicKeyFile, err := os.Create("public.pem")
	if err != nil {
		return fmt.Errorf("creating public key file: %w", err)
	}
	defer publicKeyFile.Close()

	// marshal the public key from the private key to PKIX format
	asn1Bytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		return fmt.Errorf("marshaling public key: %w", err)
	}

	// construct a pem block for the public key
	publicKeyBlock := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: asn1Bytes,
	}

	// write the pem block to the file
	err = pem.Encode(publicKeyFile, publicKeyBlock)
	if err != nil {
		return fmt.Errorf("encoding public key to file: %w", err)
	}

	fmt.Println("private and public keys generated")

	return nil
}
