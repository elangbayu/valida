package cmd

import (
	"log"

	"github.com/spf13/cobra"
	"valida/internal/apitest"
)

var file string

var testCmd = &cobra.Command{
	Use:   "test --file [JSON FILE]",
	Short: "Test the given OpenAPI Spec file",
	Long:  `Test the given OpenAPI Spec file`,
	Run: func(cmd *cobra.Command, args []string) {
		filePath, err := cmd.Flags().GetString("file")
		if err != nil {
			log.Fatal(err)
		}
		if err := apitest.TestAPISpec(filePath); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(testCmd)
	testCmd.Flags().StringVarP(&file, "file", "f", "", "OpenAPI Spec")
}
