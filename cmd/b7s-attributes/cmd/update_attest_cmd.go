package cmd

import (
	"github.com/spf13/cobra"
)

var attestCmd = &cobra.Command{
	Use:   "attest",
	Short: "Attest an existing attributes file",
	Long: `Provide an attestation to an existing attributes file. 

There are two ways to attest for an existing attributes file:
	1. Provide a private key of an attestor to sign the attributes file with.
	2. Manually specify the attestor ID and the signature`,
	RunE: runAttest,
	Args: cobra.ExactArgs(1),
}
