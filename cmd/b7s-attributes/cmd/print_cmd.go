package cmd

import (
	"github.com/spf13/cobra"
)

var flagsPrint struct {
	input string
}

// printCmd represents the create command
var printCmd = &cobra.Command{
	Use:          "print",
	Short:        "Print content of an attribute file",
	RunE:         runPrint,
	SilenceUsage: true,
}

func init() {
	printCmd.Flags().StringVarP(&flagsPrint.input, "input", "i", "", "input file to read attributes from")
}