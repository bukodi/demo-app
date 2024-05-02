package authn

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"math/big"
	"testing"
)

func TestLoadECKeyFromEnv(t *testing.T) {

	privKeyGenerated, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	privHex := hex.EncodeToString(privKeyGenerated.D.Bytes())
	t.Logf("privHex: %s", privHex)

	privKeyLoaded := new(ecdsa.PrivateKey)
	curve := elliptic.P256()
	privKeyLoaded.PublicKey.Curve = curve
	d := new(big.Int)
	dBytes, err := hex.DecodeString(privHex)
	if err != nil {
		t.Fatal(err)
	}
	d.SetBytes(dBytes)
	privKeyLoaded.D = d
	privKeyLoaded.PublicKey.X, privKeyLoaded.PublicKey.Y = curve.ScalarBaseMult(d.Bytes())

	if !privKeyLoaded.Equal(privKeyGenerated) {
		t.Fatalf("privKeyLoaded: %v, privKeyGenerated: %v", privKeyLoaded, privKeyGenerated)
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
