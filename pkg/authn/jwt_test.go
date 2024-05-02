package authn

import (
	"github.com/golang-jwt/jwt/v5"
	"testing"
	"time"
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
