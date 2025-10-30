package authn

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"math/big"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TestJWTGenerate(t *testing.T) {

	issuerKey := genEcdsaP256KeyFromSeed("my seed")

	token1 := jwt.NewWithClaims(jwt.SigningMethodES256, jwt.MapClaims{
		"foo": "bar",
		"nbf": time.Date(2015, 10, 10, 12, 0, 0, 0, time.UTC).Unix(),
	})
	tokenStr, err := token1.SignedString(issuerKey)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("tokenStr: %s", tokenStr)

	token2, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return issuerKey.Public(), nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if !token2.Valid {
		t.Fatal("token2 is invalid")
	}
	token2Claims, ok := token2.Claims.(jwt.MapClaims)
	if !ok {
		t.Fatal("token2 claims is not jwt.MapClaims")
	}
	if token2Claims["foo"] != "bar" {
		t.Fatal("token2 claims[foo] is not bar")
	}
}

func genEcdsaP256KeyFromSeed(seed string) *ecdsa.PrivateKey {
	dBytes := sha256.Sum256([]byte(seed))
	privKey := new(ecdsa.PrivateKey)
	curve := elliptic.P256()
	privKey.PublicKey.Curve = curve
	d := new(big.Int)
	d.SetBytes(dBytes[:])
	privKey.D = d
	privKey.PublicKey.X, privKey.PublicKey.Y = curve.ScalarBaseMult(d.Bytes())

	return privKey
}

func TestGenECKeyFromSeed(t *testing.T) {

	key1 := genEcdsaP256KeyFromSeed("my seed")
	msg := "Hello world"
	msgHash := sha256.Sum256([]byte(msg))
	signature, err := ecdsa.SignASN1(rand.Reader, key1, msgHash[:])
	if err != nil {
		t.Fatal(err)
	}

	key2 := genEcdsaP256KeyFromSeed("my seed")
	if !ecdsa.VerifyASN1(&key2.PublicKey, msgHash[:], signature) {
		t.Fatal("signature verification failed")
	}

	key3 := genEcdsaP256KeyFromSeed("other seed")
	if ecdsa.VerifyASN1(&key3.PublicKey, msgHash[:], signature) {
		t.Fatal("signature verification must failed")
	}
}
