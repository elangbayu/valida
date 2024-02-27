/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// execCmd represents the exec command
var execCmd = &cobra.Command{
	Use:   "exec",
	Short: "Execute automatic test for given OpenAPI Spec .json file",
	Long:  `Execute automatic test for given OpenAPI Spec .json file`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("exec called")
	},
}

func init() {
	rootCmd.AddCommand(execCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	execCmd.PersistentFlags().String("file", "", "Path to OpenAPI Spec .json file")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// execCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	execCmd.MarkPersistentFlagRequired("file")
}
