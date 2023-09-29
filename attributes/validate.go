package attributes

import (
	"encoding/base64"
	"errors"
	"fmt"
)

func Validate(att Attestation) error {

	// An unsigned attribute file cannot be attested.
	if att.Signature == nil && len(att.Attestors) > 0 {
		return errors.New("signer not found but attestor list not empty")
	}

	// If signed, verify the signature.
	if att.Signature != nil {
		data, err := serialize(att.Attributes)
		if err != nil {
			return fmt.Errorf("could not serialize attribute data: %w", err)
		}

		err = verifySignature(data, *att.Signature)
		if err != nil {
			return fmt.Errorf("could not verify signature: %w", err)
		}
	}

	// If the file has no attestations, we're done.
	if len(att.Attestors) == 0 {
		return nil
	}

	// Verify attestation signatures.
	data, err := getAttestPayload(att)
	if err != nil {
		return fmt.Errorf("could not get attestation payload: %w", err)
	}

	for i, attestor := range att.Attestors {
		err = verifySignature(data, attestor)
		if err != nil {
			return fmt.Errorf("could not verify attestor %v (%v): %w", i, attestor.Signer.String(), err)
		}
	}

	return nil
}

func verifySignature(data []byte, signature Signature) error {

	pub, err := signature.Signer.ExtractPublicKey()
	if err != nil {
		return fmt.Errorf("could not extract public key: %w", err)
	}

	sig, err := base64.StdEncoding.DecodeString(signature.Signature)
	if err != nil {
		return fmt.Errorf("could not decode signature: %w", err)
	}

	ok, err := pub.Verify(data, sig)
	if err != nil {
		return fmt.Errorf("could not verify signature: %w", err)
	}

	if !ok {
		return errors.New("signature is not valid")
	}

	return nil
}
