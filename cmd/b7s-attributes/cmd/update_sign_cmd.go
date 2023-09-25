package cmd

import (
	"github.com/spf13/cobra"
)

var signCmd = &cobra.Command{
	Use:   "sign",
	Short: "Sign an existing attributes file",
	Long: `Sign an existing attributes file. 

There are two ways to provide a signature for the attributes file:
	1. Provide a private key of a node to sign the attributes file with.
	2. Manually specify the node ID and the signature`,
	RunE: runSign,
	Args: cobra.NoArgs,
}
