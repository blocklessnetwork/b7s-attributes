package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"

	"github.com/blocklessnetwork/b7s-attributes/attributes"
	"github.com/blocklessnetwork/b7s-attributes/env"
)

var flagsCreate struct {
	prefix string
	limit  uint
	ignore []string
	strict bool
	output string
}

func runCreate(_ *cobra.Command, _ []string) error {

	attrs, err := env.ReadAttributes(flagsCreate.prefix, flagsCreate.ignore)
	if err != nil {
		return fmt.Errorf("could not retrieve attributes from the environment: %w", err)

	}

	if len(attrs) == 0 {
		log.Printf("no attributes set")
		return nil
	}

	if uint(len(attrs)) > flagsCreate.limit {

		if flagsCreate.strict {
			return fmt.Errorf("too many attributes set. remove some or re-run with 'strict' mode turned off")
		}

		attrs = attrs[:flagsCreate.limit]
	}

	log.Printf("%v attributes retrieved", len(attrs))

	att := attributes.Attestation{
		Attributes: attrs,
	}

	f, err := os.Create(flagsCreate.output)
	if err != nil {
		return fmt.Errorf("could not create output file: %w", err)
	}
	defer f.Close()

	err = attributes.Export(f, att)
	if err != nil {
		return fmt.Errorf("could not write attributes to output file: %w", err)
	}

	// TODO: Add signing.

	log.Printf("attributes written to %v", flagsCreate.output)

	return nil
}
