package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
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
