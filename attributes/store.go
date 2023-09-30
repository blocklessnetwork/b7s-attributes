package attributes

import (
	"bytes"
	"errors"
	"fmt"
	"io"

	"github.com/libp2p/go-libp2p/core/peer"
)

func ExportAttestation(writer io.Writer, att Attestation) error {

	err := att.Valid()
	if err != nil {
		return fmt.Errorf("cannot export invalid attestation: %w", err)
	}

	// Encode attribute data.
	data, err := EncodeAttributes(att.Attributes)
	if err != nil {
		return fmt.Errorf("could not encode attribute data: %w", err)
	}

	if len(data) > AttributesDataLength {
		return fmt.Errorf("attribute data too large (have: %v, limit: %v)", len(data), AttributesDataLength)
	}

	// If attribute data is too short, pad up until the size.
	if len(data) < AttributesDataLength {
		padding := make([]byte, AttributesDataLength-len(data))
		data = append(data, padding...)
	}

	// Encode signature data, if found.
	if att.Signature != nil {
		sig, err := serializeSignature(*att.Signature)
		if err != nil {
			return fmt.Errorf("could not serialize node signature: %w", err)
		}

		// Add signature data. Don't add the separator, because we'll be able to identify this data in the file by offset alone.
		data = append(data, sig...)
	}

	// Encode attestations, if any.
	for i, attestor := range att.Attestors {
		attData, err := serializeSignature(attestor)
		if err != nil {
			return fmt.Errorf("could not serialize attestor %v: %w", i, err)
		}
		// Add the separator and the signature data.
		data = append(data, binaryRecordSeparator)
		data = append(data, attData...)
	}

	_, err = writer.Write(data)
	if err != nil {
		return fmt.Errorf("could not write data to output: %w", err)
	}

	return nil
}

func ImportAttestation(reader io.Reader) (Attestation, error) {

	// First, read and decode attribute data.
	attrData := make([]byte, AttributesDataLength)
	_, err := io.ReadFull(reader, attrData)
	if err != nil {

		if !errors.Is(err, io.ErrUnexpectedEOF) {
			return Attestation{}, fmt.Errorf("could not read attestation data: %w", err)
		}

		return Attestation{}, fmt.Errorf("unexpected input data - too short attribute data payload")
	}

	attrs, err := DecodeAttributes(attrData)
	if err != nil {
		return Attestation{}, fmt.Errorf("could not decode attribute data: %w", err)
	}

	att := Attestation{
		Attributes: attrs,
	}

	signingData, err := io.ReadAll(reader)
	if err != nil {
		return Attestation{}, fmt.Errorf("could not read singing data: %w", err)
	}

	// If there's no more data, we're done.
	if len(signingData) == 0 {
		return att, nil
	}

	fields := bytes.Split(signingData, []byte{binaryRecordSeparator})
	// Fields should be in {id, sig} pairs.
	if len(fields)%2 != 0 {
		return Attestation{}, fmt.Errorf("unexpected input format, attestation data should come in pairs, have %v fields", len(fields))
	}

	// We have signing data. First try to decode node signature.
	signer, err := peer.Decode(string(fields[0]))
	if err != nil {
		return Attestation{}, fmt.Errorf("could not decode node ID (data: %s): %w", string(fields[0]), err)
	}

	sig := Signature{
		Signer:    signer,
		Signature: string(fields[1]),
	}

	att.Signature = &sig

	// If there are not more fields, we're done.
	if len(fields) == 2 {
		return att, nil
	}

	att.Attestors = make([]Signature, 0)

	for i := 2; i < len(fields); i = i + 2 {
		// We have signing data. First try to decode node signature.
		signer, err := peer.Decode(string(fields[i]))
		if err != nil {
			return Attestation{}, fmt.Errorf("could not attestor ID (no: %v, data: %s): %w", i, string(fields[i]), err)
		}

		attestor := Signature{
			Signer:    signer,
			Signature: string(fields[i+1]),
		}

		att.Attestors = append(att.Attestors, attestor)
	}

	return att, nil
}
