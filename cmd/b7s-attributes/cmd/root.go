package cmd

import (
	"log"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "b7s-attributes",
	Short: "Control panel for managing b7s node attribute file(s).",
	CompletionOptions: cobra.CompletionOptions{
		DisableDefaultCmd: true,
	},
	SilenceUsage: true,
	Args:         cobra.NoArgs,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {

	log.SetFlags(0)

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(showCmd)
	rootCmd.AddCommand(createCmd)
	rootCmd.AddCommand(updateCmd)
	rootCmd.AddCommand(printCmd)
	rootCmd.AddCommand(validateCmd)
	rootCmd.AddCommand(uploadCmd)
	rootCmd.AddCommand(addNameCmd)
}
