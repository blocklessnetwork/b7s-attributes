package cmd

import (
	"time"

	"github.com/spf13/cobra"
)

const (
	defaultValidity = 30 * 24 * time.Hour
	defaultCaching  = 2 * time.Hour
)

var addNameCmd = &cobra.Command{
	Use:          "add-name",
	Short:        "Add an IPNS name for IPFS record",
	RunE:         runAddName,
	SilenceUsage: true,
	Args:         cobra.ExactArgs(1),
}

func init() {
	addNameCmd.Flags().StringVarP(&flagsAddName.key, "key", "k", "", "key to sign the IPNS record with")
	addNameCmd.Flags().DurationVar(&flagsAddName.validity, "validity-period", defaultValidity, "how long the record will be valid for")
	addNameCmd.Flags().DurationVar(&flagsAddName.cache, "caching-period", defaultCaching, "how long should nodes cache this record")
	addNameCmd.Flags().Uint64Var(&flagsAddName.sequence, "sequence", 0, "sequence number to use for the IPNS record")
	addNameCmd.Flags().StringVar(&flagsAddName.gatewayURL, "gateway-url", "", "URL of the Blockless gateway")
}
