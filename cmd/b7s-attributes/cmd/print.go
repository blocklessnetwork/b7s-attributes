package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/blocklessnetwork/b7s-attributes/attributes"
)

func runPrint(_ *cobra.Command, args []string) error {

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

	data, err := json.Marshal(att)
	if err != nil {
		return fmt.Errorf("could not encode attributes: %w", err)
	}

	fmt.Printf("%s\n", data)

	return nil
}
