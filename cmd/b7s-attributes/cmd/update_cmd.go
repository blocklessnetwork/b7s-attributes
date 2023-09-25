package cmd

import (
	"errors"

	"github.com/spf13/cobra"
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:          "update",
	Short:        "Update an existing attributes file",
	Long:         `Update an existing attributes file by providing a signature or an attestation.`,
	SilenceUsage: true,
	Args:         cobra.NoArgs,
}

func init() {
	updateCmd.PersistentFlags().StringVarP(&flagsUpdate.input, "input", "i", "", "attribute file to update")
	updateCmd.PersistentFlags().StringVarP(&flagsUpdate.signingKey, "signing-key", "k", "", "key to sign the attribute file with")

	updateCmd.PersistentFlags().StringVar(&flagsUpdate.signerID, "signer-id", "", "signer (node) ID")
	updateCmd.PersistentFlags().StringVar(&flagsUpdate.signature, "signature", "", "signature")

	updateCmd.AddCommand(signCmd, attestCmd)
}

var flagsUpdate flagsUpdateCmd

type flagsUpdateCmd struct {
	input      string
	signingKey string
	signerID   string
	signature  string
}

func (f flagsUpdateCmd) validate() error {

	if f.input == "" {
		return errors.New("input attributes file is required")
	}

	if f.signingKey != "" {

		if f.signerID != "" || f.signature != "" {
			return errors.New("only one signing method is allowed - either a signing key or (signerID + signature) combination")
		}

		return nil
	}

	if f.signerID == "" || f.signature == "" {
		return errors.New("both signer ID and signature are required")
	}

	return nil
}
