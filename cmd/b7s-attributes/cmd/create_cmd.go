package cmd

import (
	"github.com/spf13/cobra"

	"github.com/blocklessnetwork/b7s-attributes/attributes"
)

var flagsCreate struct {
	prefix string
	limit  uint
	ignore []string
	strict bool
	output string
}

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:          "create",
	Short:        "Export attributes to a file",
	RunE:         runCreate,
	SilenceUsage: true,
}

func init() {
	createCmd.Flags().StringVar(&flagsCreate.prefix, "prefix", attributes.Prefix, "prefix node attributes environment variables have")
	createCmd.Flags().UintVarP(&flagsCreate.limit, "limit", "l", attributes.Limit, "number of node attributes to use")
	createCmd.Flags().StringSliceVarP(&flagsCreate.ignore, "ignore", "i", []string{}, "environment variables to skip")
	createCmd.Flags().BoolVar(&flagsCreate.strict, "strict", true, "stop execution if there are too many attributes")
	createCmd.Flags().StringVarP(&flagsCreate.output, "output", "o", "attributes.bin", "output file to write attributes to")
}
