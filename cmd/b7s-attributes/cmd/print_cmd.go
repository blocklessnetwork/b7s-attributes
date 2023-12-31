package cmd

import (
	"github.com/spf13/cobra"
)

// printCmd represents the create command
var printCmd = &cobra.Command{
	Use:          "print",
	Short:        "Print content of an attribute file",
	RunE:         runPrint,
	SilenceUsage: true,
	Args:         cobra.ExactArgs(1),
}
