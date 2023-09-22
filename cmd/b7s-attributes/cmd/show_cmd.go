package cmd

import (
	"github.com/spf13/cobra"

	"github.com/blocklessnetwork/b7s-attributes/attributes"
)

var showCmd = &cobra.Command{
	Use:          "show",
	Short:        "Display current node attributes",
	RunE:         runShow,
	SilenceUsage: true,
}

func init() {
	showCmd.Flags().StringVar(&flagsShow.prefix, "prefix", attributes.Prefix, "prefix node attributes environment variables have")
	showCmd.Flags().UintVarP(&flagsShow.limit, "limit", "l", attributes.Limit, "number of node attributes to use")
	showCmd.Flags().StringSliceVarP(&flagsShow.ignore, "ignore", "i", []string{}, "environment variables to skip")
	showCmd.Flags().BoolVar(&flagsShow.strict, "strict", true, "stop execution if there are too many attributes")
}
