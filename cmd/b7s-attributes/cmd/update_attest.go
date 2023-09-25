package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/spf13/cobra"

	"github.com/blocklessnetwork/b7s-attributes/attributes"
)

func runAttest(_ *cobra.Command, args []string) error {

	input := args[0]

	flags := flagsUpdate
	err := flags.validate()
	if err != nil {
		return err
	}

	if flags.signingKey != "" {
		return attestWithKey(input, flags.signingKey)
	}

	att, err := readAttributesFile(input)
	if err != nil {
		return fmt.Errorf("could not read attributes from input file: %w", err)
	}

	signerID, err := peer.Decode(flags.signerID)
	if err != nil {
		return fmt.Errorf("could not decode signer ID: %w", err)
	}

	err = addAttestation(input, att, signerID, flags.signature)
	if err != nil {
		return fmt.Errorf("could not add signature to attributes file: %w", err)
	}

	return nil
}

func attestWithKey(name string, keyPath string) error {

	key, err := readPrivateKey(keyPath)
	if err != nil {
		return fmt.Errorf("could not read key file: %w", err)
	}

	att, err := readAttributesFile(name)
	if err != nil {
		return fmt.Errorf("could not read attributes from input file: %w", err)
	}

	attestation, err := attributes.Attest(att, key)
	if err != nil {
		return fmt.Errorf("could not create attestation: %w", err)
	}

	attestorID, err := peer.IDFromPrivateKey(key)
	if err != nil {
		return fmt.Errorf("could not get attestor ID: %w", err)
	}

	return addAttestation(name, att, attestorID, attestation)
}

func addAttestation(name string, att attributes.Attestation, attestorID peer.ID, sig string) error {

	att.Attestors = append(att.Attestors, attributes.Signature{
		Signer:    attestorID,
		Signature: sig,
	})

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

	err = attributes.Export(out, att)
	if err != nil {
		return fmt.Errorf("could not write updated attributes to file: %w", err)
	}

	log.Printf("updated attributes written to %v", name)

	return nil
}
