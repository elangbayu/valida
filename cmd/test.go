package cmd

import (
	"log"
	"path/filepath"

	"valida/internal/apitest"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var file string

var testCmd = &cobra.Command{
	Use:   "test --file [JSON/YAML FILE]",
	Short: "Test the given OpenAPI Spec file",
	Long:  `Test the given OpenAPI Spec file`,
	Run: func(cmd *cobra.Command, args []string) {
		file := viper.GetString("file")
		ext := filepath.Ext(file)
		if ext != ".json" && ext != ".yaml" && ext != ".yml" {
			log.Fatal("File must be a .json, .yaml, or .yml file")
		}

		apiSpec, err := apitest.TestAPISpec(file)
		if err != nil {
			log.Fatal(err)
		}

		apitest.MakeRequest(apiSpec)
		// apitest.PrintAPISpec(apiSpec)
	},
}

func init() {
	rootCmd.AddCommand(testCmd)
	testCmd.Flags().StringVarP(&file, "file", "f", "", "OpenAPI Spec file (JSON or YAML)")
	testCmd.MarkFlagRequired("file")

	viper.BindPFlag("file", testCmd.Flags().Lookup("file"))
}
