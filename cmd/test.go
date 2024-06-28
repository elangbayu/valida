/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

// testCmd represents the test command
var testCmd = &cobra.Command{
	Use:   "test --file [JSON FILE]",
	Short: "Test the given OpenAPI Spec file",
	Long:  `Test the given OpenAPI Spec file`,
	Run: func(cmd *cobra.Command, args []string) {
		filePath, err := cmd.Flags().GetString("file")
		if err != nil {
			log.Fatal(err)
		}
		err = testAPISpec(filePath)
		if err != nil {
			log.Fatal(err)
		}
	},
}

var File string

func init() {
	rootCmd.AddCommand(testCmd)
	testCmd.Flags().StringVarP(&File, "file", "f", "", "OpenAPI Spec")
}

type Server struct {
	URL string `json:"url"`
}

type TestResult struct {
	Endpoint string
	Method   string
	Status   string
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

	var results []TestResult

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
				status := "PASSED ✅"

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

				if err != nil || resp.StatusCode >= 400 {
					status = "FAILED 🔴"
				}

				if resp != nil {
					body, errRes := io.ReadAll(resp.Body)
					errBodyClose := resp.Body.Close()
					if errBodyClose != nil {
						return errBodyClose
					}
					if errRes != nil {
						status = "FAILED 🔴"
					}

					fmt.Println("Body : ", string(body))
				}

				results = append(results, TestResult{Endpoint: endpoint, Method: strings.ToUpper(k), Status: status})
			}
		case []interface{}:
			for i, value := range v {
				fmt.Printf("Index: %d, Value: %v\n", i, value)
			}
		default:
			fmt.Println("Unknown type")
		}
	}

	displayResultsAsTable(results)
	return nil
}

func displayResultsAsTable(results []TestResult) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Endpoint", "Method", "Status"})
	// Set header to bold
	table.SetHeaderColor(
		tablewriter.Colors{tablewriter.Bold},
		tablewriter.Colors{tablewriter.Bold},
		tablewriter.Colors{tablewriter.Bold},
	)
	table.SetRowLine(true)

	for _, result := range results {
		table.Append([]string{result.Endpoint, result.Method, result.Status})
	}

	table.Render()
}
