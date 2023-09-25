package cmd

import (
	"github.com/spf13/cobra"
)

var validateCmd = &cobra.Command{
	Use:          "validate",
	Short:        "Validate attributes file",
	RunE:         runValidate,
	SilenceUsage: true,
	Args:         cobra.ExactArgs(1),
}
