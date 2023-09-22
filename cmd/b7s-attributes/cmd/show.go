package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"

	"github.com/blocklessnetwork/b7s-attributes/env"
)

var flagsShow struct {
	prefix string
	limit  uint
	ignore []string
	strict bool
}

func runShow(_ *cobra.Command, _ []string) error {

	attrs, err := env.ReadAttributes(flagsShow.prefix, flagsShow.ignore)
	if err != nil {
		return fmt.Errorf("could not retrieve attributes from the environment: %w", err)

	}

	if len(attrs) == 0 {
		log.Printf("no attributes set")
		return nil
	}

	if uint(len(attrs)) > flagsShow.limit {

		if flagsShow.strict {
			return fmt.Errorf("too many attributes set. remove some or re-run with 'strict' mode turned off")
		}

		attrs = attrs[:flagsShow.limit]
	}

	log.Printf("%v attributes retrieved", len(attrs))

	for _, attr := range attrs {
		fmt.Printf("%v: %v\n", attr.Name, attr.Value)
	}

	return nil
}
