package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var version = "dev"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version of Valida",
	Long:  `All software has versions. This is Valida's`,
	Run: func(cmd *cobra.Command, args []string) {
		printVersion()
	},
}

func printVersion() {
	fmt.Printf("Valida %s\n", version)
}
