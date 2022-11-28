package main

import (
	"context"
	"fmt"
	"time"

	"github.com/dtherhtun/service/business/data/schema"
	"github.com/dtherhtun/service/business/sys/database"
)

func main() {
	err := migrate()
	if err != nil {
		fmt.Println(err)
	}
}

func seed() error {
	cfg := database.Config{
		User:         "postgres",
		Password:     "postgres",
		Host:         "localhost",
		Name:         "postgres",
		MaxIdleConns: 0,
		MaxOpenConns: 0,
		DisableTLS:   true,
	}
	db, err := database.Open(cfg)
	if err != nil {
		return fmt.Errorf("connect database: %w", err)
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	if err := schema.Seed(ctx, db); err != nil {
		return fmt.Errorf("seed database: %w", err)
	}

	fmt.Println("seed data complete")
	return nil
}

func migrate() error {
	cfg := database.Config{
		User:         "postgres",
		Password:     "postgres",
		Host:         "localhost",
		Name:         "postgres",
		MaxIdleConns: 0,
		MaxOpenConns: 0,
		DisableTLS:   true,
	}
	db, err := database.Open(cfg)
	if err != nil {
		return fmt.Errorf("connect database: %w", err)
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := schema.Migrate(ctx, db); err != nil {
		return fmt.Errorf("migrate database: %w", err)
	}

	fmt.Println("migration complete")
	return seed()
}

//func genToken() error {
//
//	name := "zarf/keys/54bb2165-71e1-41a6-af3e-7da4a0e1e2c1.pem"
//	f, err := os.Open(name)
//	if err != nil {
//		return err
//	}
//	defer f.Close()
//
//	// limit PEM file size to 1 megabyte. This should be reasonable for
//	// almost any PEM file and prevents shenanigans like linking the file
//	// to /dev/random or something like that.
//	privatePEM, err := io.ReadAll(io.LimitReader(f, 1024*1024))
//	if err != nil {
//		return fmt.Errorf("reading auth private key: %w", err)
//	}
//
//	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privatePEM)
//	if err != nil {
//		return fmt.Errorf("parsing auth private key: %w", err)
//	}
//
//	// =======================================================================
//
//	// Generating a token requires defining a set of the claims. In this applications
//	// case, we only care about defining the subject and the user in question and
//	// the roles they have on the database. This token will expire in a year.
//	//
//	// iss (issuer): Issuer of the JWT
//	// sub (subject): Subject of the JWT (the user)
//	// aud (audience): Recipient for which the JWT is intended
//	// exp (expiration time): Time after which the JWT expires
//	// nbf (not before time): Time before which the JWT must not be accepted for processing
//	// iat (issued at time): Time at which the JWT was issued; can be used to determine age of the JWT
//	// jti (JWT ID): Unique identifier; can be used to prevent the JWT from being replayed (allows a token to be used only one time)
//	claims := struct {
//		jwt.RegisteredClaims
//		Roles []string
//	}{
//		RegisteredClaims: jwt.RegisteredClaims{
//			Issuer:    "service project",
//			Subject:   "123456789",
//			ExpiresAt: &jwt.NumericDate{Time: time.Now().Add(8760 * time.Hour)},
//			IssuedAt:  &jwt.NumericDate{Time: time.Now().UTC()},
//		},
//		Roles: []string{"ADMIN"},
//	}
//
//	method := jwt.GetSigningMethod("RS256")
//	token := jwt.NewWithClaims(method, claims)
//	token.Header["kid"] = "54bb2165-71e1-41a6-af3e-7da4a0e1e2c1"
//
//	tokenStr, err := token.SignedString(privateKey)
//	if err != nil {
//		return err
//	}
//
//	fmt.Println("========== TOKEN BEGIN ==========")
//	fmt.Println(tokenStr)
//	fmt.Println("========== TOKEN END ==========")
//	fmt.Print("\n")
//
//	// =====================================================================
//
//	// Marshal the public key from the private key to PKIX
//	asn1Bytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
//	if err != nil {
//		return fmt.Errorf("marshalling public key: %w", err)
//	}
//
//	publicBlock := pem.Block{
//		Type:  "RSA PUBLIC KEY",
//		Bytes: asn1Bytes,
//	}
//
//	if err := pem.Encode(os.Stdout, &publicBlock); err != nil {
//		return fmt.Errorf("encoding to public file: %w", err)
//	}
//
//	// ======================================================================
//
//	fmt.Println("=======================")
//	// Create the token parser to use. The algorithm used to sign the JWT must be
//	// validated to avoid a critical vulnerability:
//	// http://auth0.com/blog/critical-vulnerabilities-in-json-web-token-libraries
//	//parser := jwt.Parser{
//	//	ValidMethods: []string{"RS256"},
//	//}
//
//	option := jwt.WithValidMethods([]string{"RS256"})
//
//	parser := jwt.NewParser(option)
//
//	keyFunc := func(t *jwt.Token) (interface{}, error) {
//		kid, ok := t.Header["kid"]
//		if !ok {
//			return nil, errors.New("missing key id (kid) in token header")
//		}
//		kidID, ok := kid.(string)
//		if !ok {
//			return nil, errors.New("user token key id (kid) must be string")
//		}
//		fmt.Println("KID:", kidID)
//
//		return &privateKey.PublicKey, nil
//	}
//
//	var parsedClaims struct {
//		jwt.RegisteredClaims
//		Roles []string
//	}
//
//	parsedToken, err := parser.ParseWithClaims(tokenStr, &parsedClaims, keyFunc)
//	if err != nil {
//		return fmt.Errorf("parsing token: %w", err)
//	}
//
//	if !parsedToken.Valid {
//		return errors.New("invalid token")
//	}
//
//	fmt.Println("Token Validate")
//
//	return nil
//}
//
//func genKey() error {
//
//	// Generate new private key.
//	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
//	if err != nil {
//		return err
//	}
//
//	privateFile, err := os.Create("private.pem")
//	if err != nil {
//		return fmt.Errorf("creating private file: %w", err)
//	}
//	defer privateFile.Close()
//
//	// Construct a PEM block for the private key.
//	privateBlock := pem.Block{
//		Type:  "RSA PRIVATE KEY",
//		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
//	}
//
//	if err := pem.Encode(privateFile, &privateBlock); err != nil {
//		return fmt.Errorf("encoding to private file: %w", err)
//	}
//
//	// Marshal the public key from the private key to PKIX
//	asn1Bytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
//	if err != nil {
//		return fmt.Errorf("marshalling public key: %w", err)
//	}
//
//	publicFile, err := os.Create("public.pem")
//	if err != nil {
//		return fmt.Errorf("creating public file: %w", err)
//	}
//	defer publicFile.Close()
//
//	publicBlock := pem.Block{
//		Type:  "RSA PUBLIC KEY",
//		Bytes: asn1Bytes,
//	}
//
//	if err := pem.Encode(publicFile, &publicBlock); err != nil {
//		return fmt.Errorf("encoding to public file: %w", err)
//	}
//
//	fmt.Println("private and public key files generated")
//
//	return nil
//}
