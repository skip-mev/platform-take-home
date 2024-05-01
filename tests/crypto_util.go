package tests

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/sha256"
	"crypto/x509"
	"fmt"
	"math/big"
)

func pubKeyFromBytes(bz []byte) *ecdsa.PublicKey {
	pubKeyInterface, err := x509.ParsePKIXPublicKey(bz)
	if err != nil {
		fmt.Println("Error parsing public key:", err)
		return nil
	}

	pk, ok := pubKeyInterface.(*ecdsa.PublicKey)
	if !ok {
		fmt.Println("Not an ECDSA public key")
		return nil
	}

	return pk
}

func signatureFromBytes(sigStr []byte) *signature {
	return &signature{
		R: new(big.Int).SetBytes(sigStr[:32]),
		S: new(big.Int).SetBytes(sigStr[32:64]),
	}
}

var p256Order = elliptic.P256().Params().N

var p256HalfOrder = new(big.Int).Rsh(p256Order, 1)

// signature holds the r and s values of an ECDSA signature.
type signature struct {
	R, S *big.Int
}

func VerifySignature(pk ecdsa.PublicKey, msg, sig []byte) bool {
	// check length for raw signature
	// which is two 32-byte padded big.Ints
	// concatenated
	// NOT DER!

	if len(sig) != 64 {
		return false
	}

	s := signatureFromBytes(sig)
	if !IsSNormalized(s.S) {
		return false
	}

	h := sha256.Sum256(msg)
	return ecdsa.Verify(&pk, h[:], s.R, s.S)
}

func IsSNormalized(sigS *big.Int) bool {
	return sigS.Cmp(p256HalfOrder) != 1
}
