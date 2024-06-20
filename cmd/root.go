/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "valida --test spec.json",
	Short: "Automatic API Testing Execution",
	Long:  `Run API automation testing with your OpenAPI Spec .json file`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		filePath, err := cmd.Flags().GetString("test")
		if err != nil {
			log.Fatal(err)
		}
		err = testAPISpec(filePath)
		if err != nil {
			log.Fatal(err)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	//cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	//rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.validarc.json)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.PersistentFlags().String("test", "", "path to the OpenAPI specification JSON file (required)")
	rootCmd.MarkPersistentFlagRequired("test")
}

type Server struct {
	URL string `json:"url"`
}

func testAPISpec(filePath string) error {
	spec, err := openapi3.NewLoader().LoadFromFile(filePath)
	if err != nil {
		return err
	}

	// Implement the API testing logic here
	fmt.Println("Successfully parsed OpenAPI Specification:")
	fmt.Println("Title:", spec.Info.Title)
	fmt.Println("Version:", spec.Info.Version)
	// serve, _ := spec.Servers.BasePath();
	// fmt.Println("Server:", string(spec.Servers))
	fmt.Println("Testing APIs... (this is a stub)")

	data, err := json.MarshalIndent(spec.Paths, "", "  ")
	if err != nil {
		fmt.Println("error:", err)
	}

	// fmt.Println(string(data))

	// Deklarasikan map untuk menampung hasil parsing JSON
	var paths map[string]interface{}

	// Parsing JSON
	if err := json.Unmarshal([]byte(data), &paths); err != nil {
		log.Fatalf("Error parsing JSON: %v", err)
	}

	dataServers, err := json.MarshalIndent(spec.Servers, "", "  ")
	if err != nil {
		fmt.Println("error:", err)
	}

	// fmt.Println(string(dataServers))
	// var server spec.Servers

	// Unmarshal dataServers ke dalam struktur Go
	var servers []Server
	err = json.Unmarshal(dataServers, &servers)
	if err != nil {
		log.Fatalf("Error unmarshalling JSON: %v", err)
	}

	// Ambil nilai dari parameter "url"
	if len(servers) > 0 {
		url := servers[0].URL
		// fmt.Println("URL:", url)
		// Iterasi melalui map dan cetak setiap key yang merupakan endpoint
		for endpoint := range paths {
			fmt.Println("Yang dites: ", url, endpoint)
			resp, err := http.Get(url + endpoint)
			if err != nil {
				fmt.Println("Error: ", err)
			}
			defer resp.Body.Close()
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				fmt.Println("Error: ", err)
			}
			fmt.Println(string(body))
		}
	} else {
		fmt.Println("No servers found")
	}
	return nil
}
