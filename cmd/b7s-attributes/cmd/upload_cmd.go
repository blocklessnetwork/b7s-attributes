package cmd

import (
	"github.com/spf13/cobra"
)

var uploadCmd = &cobra.Command{
	Use:          "upload",
	Short:        "Upload attributes file to IPFS",
	RunE:         runUpload,
	SilenceUsage: true,
	Args:         cobra.ExactArgs(1),
}

func init() {
	uploadCmd.Flags().StringVar(&flagsUpload.gatewayURL, "gateway-url", "", "URL of the Blockless gateway")
}
