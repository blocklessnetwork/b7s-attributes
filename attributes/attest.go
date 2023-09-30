package attributes

import (
	"encoding/base64"
	"errors"
	"fmt"

	"github.com/libp2p/go-libp2p/core/crypto"
)

func Attest(att Attestation, key crypto.PrivKey) (string, error) {

	if att.Signature == nil {
		return "", errors.New("attributes have to be signed first")
	}

	data, err := getAttestPayload(att)
	if err != nil {
		return "", fmt.Errorf("could not get payload for attestor signing: %w", err)
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

func getAttestPayload(att Attestation) ([]byte, error) {

	if att.Signature == nil {
		return nil, ErrNoSignature
	}

	data, err := serializeAttributes(att.Attributes)
	if err != nil {
		return nil, fmt.Errorf("could not serialize attribute data: %w", err)
	}

	sig, err := serializeSignature(*att.Signature)
	if err != nil {
		return nil, fmt.Errorf("could not serialize node signature: %w", err)
	}

	data = append(data, sig...)

	return data, nil
}
