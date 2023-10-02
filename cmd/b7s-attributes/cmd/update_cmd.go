package cmd

import (
	"github.com/spf13/cobra"
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:          "update",
	Short:        "Update an existing attributes file",
	Long:         `Update an existing attributes file by providing a signature or an attestation.`,
	SilenceUsage: true,
	Args:         cobra.ExactArgs(1),
}

func init() {
	updateCmd.PersistentFlags().StringVarP(&flagsUpdate.signingKey, "signing-key", "k", "", "key to sign the attribute file with")

	updateCmd.AddCommand(signCmd, attestCmd)
}

var flagsUpdate flagsUpdateCmd

type flagsUpdateCmd struct {
	signingKey string
}
