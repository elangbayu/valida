/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/spf13/cobra"
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

var file string

func init() {
	rootCmd.AddCommand(testCmd)
	testCmd.Flags().StringVarP(&file, "file", "f", "", "OpenAPI Spec")
}

type Server struct {
	URL string `json:"url"`
}

type TestResult struct {
	Endpoint string
	Method   string
	Status   string
}

var httpClient = &http.Client{
	Timeout: 10 * time.Second,
}

func makeRequest(ctx context.Context, method, url string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("making request: %w", err)
	}
	return resp, nil
}

func testAPISpec(filePath string) error {
	spec, err := openapi3.NewLoader().LoadFromFile(filePath)
	if err != nil {
		return fmt.Errorf("loading OpenAPI spec: %w", err)
	}

	// Validate File OpenApi Specification
	loader := openapi3.NewLoader()
	if err := spec.Validate(loader.Context); err != nil {
		log.Fatalf("Error validating OpenAPI spec: %v", err)
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
		return fmt.Errorf("unmarshalling servers: %w", err)
	}

	if len(servers) == 0 {
		return fmt.Errorf("no servers found in the specification")
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
				status := "PASSED"

				ctx := context.Background()

				switch k {
				case "get":
					resp, err = makeRequest(ctx, http.MethodGet, fullURL)
				case "post":
					resp, err = makeRequest(ctx, http.MethodPost, fullURL)
				case "put":
					resp, err = makeRequest(ctx, http.MethodPut, fullURL)
				case "delete":
					resp, err = makeRequest(ctx, http.MethodDelete, fullURL)
				case "patch":
					resp, err = makeRequest(ctx, http.MethodPatch, fullURL)
				default:
					return fmt.Errorf("unsupported http method")
				}

				if err != nil || resp.StatusCode >= 400 {
					status = "FAILED"
				}

				if resp != nil {
					body, errRes := io.ReadAll(resp.Body)
					defer resp.Body.Close()
					if errRes != nil {
						status = "FAILED"
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

	// Sort results by Endpoint and Method
	sort.Slice(results, func(i, j int) bool {
		if results[i].Endpoint == results[j].Endpoint {
			return results[i].Method < results[j].Method
		}
		return results[i].Endpoint < results[j].Endpoint
	})

	displayResultsAsTable(results)
	return nil
}

func displayResultsAsTable(results []TestResult) {
	const (
		green = lipgloss.Color("#009900")
		red   = lipgloss.Color("#f61901")
	)

	re := lipgloss.NewRenderer(os.Stdout)

	var (
		HeaderStyle       = re.NewStyle().Foreground(green).Bold(true).Align(lipgloss.Center)
		CellStyle         = re.NewStyle().Padding(0, 1).Width(14)
		PassedStatusStyle = re.NewStyle().Foreground(green).Width(14)
		FailedStatusStyle = re.NewStyle().Foreground(red).Width(14)
		BorderStyle       = lipgloss.NewStyle().Foreground(green)
	)

	t := table.New().
		Border(lipgloss.ThickBorder()).
		BorderStyle(BorderStyle).
		Headers("ENDPOINT", "METHOD", "STATUS").
		StyleFunc(func(row, col int) lipgloss.Style {
			if row == 0 {
				return HeaderStyle
			}

			// Check if the row is within the range of results
			if row-1 < len(results) {
				// Apply styles for the status column
				if col == 2 { // The status column is the third column (index 2)
					switch results[row-1].Status {
					case "PASSED":
						return PassedStatusStyle
					case "FAILED":
						return FailedStatusStyle
					}
				}
			}

			return CellStyle
		})

	// Counters for passed and failed statuses
	var passedCount, failedCount int

	// Iterate over the testResults slice and convert each TestResult to a slice of strings
	for _, result := range results {
		t.Row(result.Endpoint, result.Method, result.Status)
		if result.Status == "PASSED" {
			passedCount++
		} else if result.Status == "FAILED" {
			failedCount++
		}
	}

	// Add rows for total passed and failed endpoints
	t.Row("", "", "") // Empty row for spacing
	t.Row("PASSED", fmt.Sprintf("%d", passedCount), "")
	t.Row("FAILED", fmt.Sprintf("%d", failedCount), "")

	fmt.Println(t)
}
