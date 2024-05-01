package server

import (
	"encoding/asn1"
	"math/big"
)

type ECDSASignature struct {
	R, S *big.Int
}

func DecodeDERSignature(der []byte) (*ECDSASignature, error) {
	var sig ECDSASignature
	_, err := asn1.Unmarshal(der, &sig)
	if err != nil {
		return nil, err
	}
	return &sig, nil
}
