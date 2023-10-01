package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"

	"github.com/blocklessnetwork/b7s-attributes/attributes"
	"github.com/blocklessnetwork/b7s-attributes/w3s"
)

func runUpload(_ *cobra.Command, args []string) error {

	input := args[0]

	f, err := os.Open(input)
	if err != nil {
		return fmt.Errorf("could not open attributes file: %w", err)
	}
	defer f.Close()

	// Make sure the file is vaid before uploading it.
	att, err := attributes.ImportAttestation(f)
	if err != nil {
		return fmt.Errorf("could not read attribute file: %w", err)
	}

	err = attributes.Validate(att)
	if err != nil {
		return fmt.Errorf("attribute file is invalid: %w", err)
	}

	// Seek back to start of file to upload it (since we've already read it).
	_, err = f.Seek(0, io.SeekStart)
	if err != nil {
		return fmt.Errorf("could not seek back to start of attributes file: %w", err)
	}

	cid, err := w3s.Upload(f)
	if err != nil {
		return fmt.Errorf("could not upload attributes file: %w", err)
	}

	fmt.Printf("%v\n", cid.String())

	return nil
}
