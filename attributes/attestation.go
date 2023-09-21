package attributes

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/libp2p/go-libp2p/core/peer"
)

const (
	AttributesDataLength = 1024 // 1KiB
)

type Signature struct {
	Signer    peer.ID
	Signature string
}

type Attestation struct {
	Attributes []Attribute
	Signature  *Signature
	Attestors  []Signature
}

func (n Attestation) Valid() error {

	if len(n.Attributes) == 0 {
		return errors.New("no attributes found")
	}

	if n.Signature == nil && len(n.Attestors) > 0 {
		return errors.New("node signature must be set before attestations")
	}

	return nil
}

func (n *Signature) serialize() ([]byte, error) {

	var buf bytes.Buffer
	_, err := buf.WriteString(fmt.Sprintf("%s%c%s%c", n.Signer.String(), binaryRecordSeparator, n.Signature, binaryRecordSeparator))
	if err != nil {
		return nil, fmt.Errorf("could not write attribute data: %w", err)
	}

	return buf.Bytes(), nil
}
