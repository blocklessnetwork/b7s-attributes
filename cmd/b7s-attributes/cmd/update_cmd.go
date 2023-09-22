package cmd

import (
	"github.com/spf13/cobra"
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update an existing attributes file",
	Long: `Update an existing attributes file. 

It is possible to provide a private key of a node to sign the attributes file.
On the other hand, it's possible to manually specify the node ID and the signature 
generated manually, for example using the b7s-keyforge utility.`,
	RunE:         runUpdate,
	SilenceUsage: true,
}

// TODO: Think of using subcommands to handle signing attributes (node) vs adding attestation (3rd party attestor).

func init() {
	updateCmd.Flags().StringVarP(&flagsUpdate.input, "input", "i", "", "attribute file to update")
	updateCmd.Flags().StringVarP(&flagsUpdate.signingKey, "key", "k", "", "key to sign the attribute file with")

	updateCmd.Flags().StringVar(&flagsUpdate.signerID, "signer-id", "", "signer (node) ID")
	updateCmd.Flags().StringVar(&flagsUpdate.signature, "signature", "", "signature")
}
