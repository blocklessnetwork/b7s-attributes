package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/blocklessnetwork/b7s-attributes/attributes"
)

func runValidate(_ *cobra.Command, args []string) error {

	input := args[0]

	f, err := os.Open(input)
	if err != nil {
		return fmt.Errorf("could not open attribute file: %w", err)
	}
	defer f.Close()

	att, err := attributes.ImportAttestation(f)
	if err != nil {
		return fmt.Errorf("could not read attribute file: %w", err)
	}

	err = attributes.Validate(att)
	if err != nil {
		return fmt.Errorf("attribute file is invalid: %w", err)
	}

	fmt.Printf("OK\n")

	return nil
}
