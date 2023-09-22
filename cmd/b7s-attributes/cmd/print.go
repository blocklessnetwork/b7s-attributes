package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/blocklessnetwork/b7s-attributes/attributes"
	"github.com/spf13/cobra"
)

func runPrint(_ *cobra.Command, _ []string) error {

	if flagsPrint.input == "" {
		return errors.New("input file is required")
	}

	f, err := os.Open(flagsPrint.input)
	if err != nil {
		return fmt.Errorf("could not open attribute file: %w", err)
	}
	defer f.Close()

	att, err := attributes.Import(f)
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
