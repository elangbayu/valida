/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/spf13/cobra"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
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

	fmt.Println("Successfully parsed OpenAPI Specification:")
	fmt.Println("Title:", spec.Info.Title)
	fmt.Println("Version:", spec.Info.Version)
	fmt.Println("Testing APIs... (this is a stub)")

	dataServers, err := json.MarshalIndent(spec.Servers, "", "  ")
	if err != nil {
		fmt.Println("error:", err)
	}

	var servers []Server
	err = json.Unmarshal(dataServers, &servers)
	if err != nil {
		log.Fatalf("Error unmarshalling JSON: %v", err)
	}

	if len(servers) == 0 {
		fmt.Println("No servers found")
		return nil
	}

	url := servers[0].URL
	fmt.Println("URL Server: ", url)

	data, err := json.MarshalIndent(spec.Paths, "", "  ")
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Println(string(data))

	var paths map[string]interface{}
	// Parsing JSON
	if err := json.Unmarshal(data, &paths); err != nil {
		log.Fatalf("Error parsing JSON: %v", err)
	}

	for key, val := range paths {
		// key is an endpoint
		fmt.Println(key)
		endpoint := strings.Replace(key, "{resource}", "test", -1) // Example: Replace `{resource}` with `users`
		endpoint = strings.Replace(endpoint, "{id}", "1", -1)      // Example: Replace `{id}` with `1`
		switch v := val.(type) {
		case map[string]interface{}:
			for k := range v {
				// k is a method
				fullURL := url + endpoint
				fmt.Println("Full URL : ", fullURL)
				fmt.Println("Method : ", k)
				var resp *http.Response
				var err error
				if k == "get" {
					resp, err = http.Get(fullURL)
				} else if k == "post" {
					resp, err = http.Post(fullURL, "application/json", nil)
				} else if k == "put" {
					req, err := http.NewRequest(http.MethodPut, fullURL, nil)
					if err == nil {
						client := &http.Client{}
						resp, err = client.Do(req)
					}
				} else if k == "patch" {
					req, err := http.NewRequest(http.MethodPatch, fullURL, nil)
					if err == nil {
						client := &http.Client{}
						resp, err = client.Do(req)
					}
				} else if k == "delete" {
					req, err := http.NewRequest(http.MethodDelete, fullURL, nil)
					if err == nil {
						client := &http.Client{}
						resp, err = client.Do(req)
					}
				}

				if err != nil {
					fmt.Printf("Error: %v\n", err)
					continue
				}

				body, errRes := io.ReadAll(resp.Body)
				errBodyClose := resp.Body.Close()
				if errBodyClose != nil {
					return errBodyClose
				}
				if errRes != nil {
					fmt.Printf("Error: %v\n", err)
					continue
				}

				fmt.Println("Body : ", string(body))

			}
		case []interface{}:
			for i, value := range v {
				fmt.Printf("Index: %d, Value: %v\n", i, value)
			}
		default:
			fmt.Println("Unknown type")
		}
	}

	return nil
}
