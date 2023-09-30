package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/spf13/cobra"

	"github.com/blocklessnetwork/b7s-attributes/attributes"
)

func runSign(_ *cobra.Command, args []string) error {

	input := args[0]

	flags := flagsUpdate
	err := flags.validate()
	if err != nil {
		return err
	}

	if flags.signingKey != "" {
		return signWithKey(input, flags.signingKey)
	}

	att, err := readAttributesFile(input)
	if err != nil {
		return fmt.Errorf("could not read attributes from input file: %w", err)
	}

	signerID, err := peer.Decode(flags.signerID)
	if err != nil {
		return fmt.Errorf("could not decode signer ID: %w", err)
	}

	err = addSignature(input, att, signerID, flags.signature)
	if err != nil {
		return fmt.Errorf("could not add signature to attributes file: %w", err)
	}

	return nil
}

func signWithKey(input string, keyPath string) error {

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

	err = addSignature(input, att, signerID, signature)
	if err != nil {
		return fmt.Errorf("could not add signature to attributes file: %w", err)
	}

	return nil
}

func addSignature(name string, att attributes.Attestation, signerID peer.ID, signature string) error {

	att.Signature = &attributes.Signature{
		Signer:    signerID,
		Signature: signature,
	}

	oldfile := name + ".old"
	err := os.Rename(name, oldfile)
	if err != nil {
		return fmt.Errorf("could not backup original attributes file")
	}

	log.Printf("old attributes file moved to %v", oldfile)

	out, err := os.Create(name)
	if err != nil {
		return fmt.Errorf("could not open file for writing update attributes: %w", err)
	}
	defer out.Close()

	err = attributes.ExportAttestation(out, att)
	if err != nil {
		return fmt.Errorf("could not write updated attributes to file: %w", err)
	}

	log.Printf("updated attributes written to %v", name)

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

	att, err := attributes.ImportAttestation(f)
	if err != nil {
		return attributes.Attestation{}, fmt.Errorf("could not import attributes from a file: %w", err)
	}

	return att, nil
}
