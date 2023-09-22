package cmd

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/blocklessnetwork/b7s-attributes/attributes"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/spf13/cobra"
)

var flagsUpdate struct {
	input      string
	signingKey string

	signerID  string
	signature string
}

func runUpdate(_ *cobra.Command, _ []string) error {

	flags := flagsUpdate
	if flags.input == "" {
		return errors.New("input attributes file is required")
	}

	if flags.signingKey == "" && flags.signerID == "" && flags.signature == "" {
		return errors.New("either a signing key, or (signerID + signature) combination is required")
	}

	// We only support one of two options - either provide the signing key, or ID+sig.
	if flags.signingKey != "" {

		if flags.signerID != "" || flags.signature != "" {
			return errors.New("only one signing method is allowed - either a signing key or (signerID + signature) combination")
		}

		return runUpdateSigningKey(flags.input, flags.signingKey)
	}

	if flags.signerID == "" || flags.signature == "" {
		return errors.New("both signer ID and signature are required")
	}

	return fmt.Errorf("TBD: not implemented")
}

func runUpdateSigningKey(input string, keyPath string) error {

	key, err := readPrivateKey(keyPath)
	if err != nil {
		return fmt.Errorf("could not read key file: %w", err)
	}

	att, err := readAttributesFile(input)
	if err != nil {
		return fmt.Errorf("could not read attributes from input file: %w", err)
	}

	signature, err := attributes.SignAttributes(att.Attributes, key)
	if err != nil {
		return fmt.Errorf("could not sign attribute data: %w", err)
	}

	signerID, err := peer.IDFromPrivateKey(key)
	if err != nil {
		return fmt.Errorf("could not get signer ID: %w", err)
	}

	att.Signature = &attributes.Signature{
		Signer:    signerID,
		Signature: signature,
	}

	oldfile := input + ".old"
	err = os.Rename(input, oldfile)
	if err != nil {
		return fmt.Errorf("could not backup original attributes file")
	}

	log.Printf("old attributes file moved to %v", oldfile)

	out, err := os.Create(input)
	if err != nil {
		return fmt.Errorf("could not open file for writing update attributes: %w", err)
	}
	defer out.Close()

	err = attributes.Export(out, att)
	if err != nil {
		return fmt.Errorf("could not write updated attributes to file: %w", err)
	}

	log.Printf("updated attributes written to %v", input)

	return nil
}

func readPrivateKey(filepath string) (crypto.PrivKey, error) {

	payload, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("could not read file: %w", err)
	}

	key, err := crypto.UnmarshalPrivateKey(payload)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal private key: %w", err)
	}

	return key, nil
}

func readAttributesFile(input string) (attributes.Attestation, error) {

	f, err := os.Open(input)
	if err != nil {
		return attributes.Attestation{}, fmt.Errorf("could not open attributes file: %w", err)
	}
	defer f.Close()

	att, err := attributes.Import(f)
	if err != nil {
		return attributes.Attestation{}, fmt.Errorf("could not import attributes from a file: %w", err)
	}

	return att, nil
}
