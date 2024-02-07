package main

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	_ "embed"
	"encoding/pem"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/open-policy-agent/opa/rego"
)

var (
	//go:embed rego/authenticate.rego
	opaAuthentication string
)

func main() {

	err := gentoken()

	if err != nil {
		log.Fatalln(err)
	}

}

func gentoken() error {

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return fmt.Errorf("generating private key: %w", err)
	}

	// Generating a token requires defining a set of claims. In this applications
	// case, we only care about defining the subject and the user in question and
	// the roles they have on the database. This token will expire in a year.
	//
	// iss (issuer): Issuer of the JWT
	// sub (subject): Subject of the JWT (the user)
	// aud (audience): Recipient for which the JWT is intended
	// exp (expiration time): Time after which the JWT expires
	// nbf (not before time): Time before which the JWT must not be accepted for processing
	// iat (issued at time): Time at which the JWT was issued; can be used to determine age of the JWT
	// jti (JWT ID): Unique identifier; can be used to prevent the JWT from being replayed (allows a token to be used only once)
	claims := struct {
		jwt.RegisteredClaims
		Roles []string
	}{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   "12345678789",
			Issuer:    "service project",
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(8760 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		},
		Roles: []string{"ADMIN"},
	}

	method := jwt.GetSigningMethod(jwt.SigningMethodRS256.Name)
	token := jwt.NewWithClaims(method, claims)
	token.Header["kid"] = "54bb2165-71e1-41a6-af3e-7da4a0e1e2c1"

	str, err := token.SignedString(privateKey)
	if err != nil {
		return fmt.Errorf("signing token: %w", err)
	}

	fmt.Println("***********************")
	fmt.Println(str)
	fmt.Println("***********************")

	// ============================================================================================================

	var claims2 struct {
		jwt.RegisteredClaims
		Roles []string
	}
	parser := jwt.NewParser(jwt.WithValidMethods([]string{jwt.SigningMethodRS256.Name}))
	parsedToken, err := parser.ParseWithClaims(str, &claims2, func(t *jwt.Token) (interface{}, error) {
		return &privateKey.PublicKey, nil
	})
	if err != nil {
		return fmt.Errorf("parsing token: %w", err)
	}
	if !parsedToken.Valid {
		return errors.New("token is invalid")
	}

	fmt.Println("SIGNATURE VERIFIED")
	fmt.Printf("claims: %#v\n", claims2)
	fmt.Println("***********************")

	// ============================================================================================================
	var claims3 struct {
		jwt.RegisteredClaims
		Roles []string
	}
	_, _, err = parser.ParseUnverified(str, &claims3)
	if err != nil {
		return fmt.Errorf("parsing token unverified: %w", err)
	}

	// Marshal the public key from the private key to PKIX.
	asn1Bytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		return fmt.Errorf("marshaling public key: %w", err)
	}

	// Construct a PEM block for the public key.
	publicBlock := pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: asn1Bytes,
	}

	var b bytes.Buffer

	// Write the public key to the public key file.
	if err := pem.Encode(&b, &publicBlock); err != nil {
		return fmt.Errorf("encoding to public file: %w", err)
	}

	input := map[string]any{
		"Key":   b.String(),
		"Token": str,
	}

	if err := opaPolicyEvaluation(context.Background(), opaAuthentication, input); err != nil {
		return fmt.Errorf("authentication failed : %w", err)
	}

	fmt.Println("SIGNATURE VALIDATED BY REGO")
	fmt.Println("****************")

	return nil
}

func opaPolicyEvaluation(ctx context.Context, opaPolicy string, input any) error {
	const opaPackage = "ardan.rego"
	const rule string = "auth"

	query := fmt.Sprintf("x = data.%s.%s", opaPackage, rule)

	q, err := rego.New(
		rego.Query(query),
		rego.Module("policy.rego", opaPolicy),
	).PrepareForEval(ctx)
	if err != nil {
		return err
	}

	results, err := q.Eval(ctx, rego.EvalInput(input))
	if err != nil {
		return fmt.Errorf("query: %w", err)
	}

	if len(results) == 0 {
		return errors.New("no results")
	}

	result, ok := results[0].Bindings["x"].(bool)
	if !ok || !result {
		return fmt.Errorf("bindings results[%v] ok[%v]", results, ok)
	}

	return nil
}

func genkey() (*rsa.PrivateKey, error) {

	// Generate a new key pair
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, fmt.Errorf("generating private key: %w", err)
	}

	// create a file for the private key in pem format
	privateKeyFile, err := os.Create("private.pem")
	if err != nil {
		return nil, fmt.Errorf("creating private key file: %w", err)
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
		return nil, fmt.Errorf("encoding private key to file: %w", err)
	}

	// create a file for the public key in pem format
	publicKeyFile, err := os.Create("public.pem")
	if err != nil {
		return nil, fmt.Errorf("creating public key file: %w", err)
	}
	defer publicKeyFile.Close()

	// marshal the public key from the private key to PKIX format
	asn1Bytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		return nil, fmt.Errorf("marshaling public key: %w", err)
	}

	// construct a pem block for the public key
	publicKeyBlock := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: asn1Bytes,
	}

	// write the pem block to the file
	err = pem.Encode(publicKeyFile, publicKeyBlock)
	if err != nil {
		return nil, fmt.Errorf("encoding public key to file: %w", err)
	}

	fmt.Println("private and public keys generated")

	return privateKey, nil
}
