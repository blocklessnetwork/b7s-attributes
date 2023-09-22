package attributes

import (
	"encoding/base64"
	"fmt"

	"github.com/libp2p/go-libp2p/core/crypto"
)

func SignAttributes(attrs []Attribute, key crypto.PrivKey) (string, error) {

	data, err := serialize(attrs)
	if err != nil {
		return "", fmt.Errorf("could not serialize attribute data: %w", err)
	}

	sig, err := key.Sign(data)
	if err != nil {
		return "", fmt.Errorf("could not sign attribute data: %w", err)
	}

	encodedSig := base64.StdEncoding.EncodeToString(sig)
	if err != nil {
		return "", fmt.Errorf("coult not encode signature: %w", err)
	}

	return encodedSig, nil
}
