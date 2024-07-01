package apitest

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
)

func TestAPISpec(filePath string) error {
	spec, err := loadAndValidateSpec(filePath)
	if err != nil {
		return err
	}

	printSpecInfo(spec)

	url, err := getServerURL(spec)
	if err != nil {
		return err
	}

	paths, err := getPathsFromSpec(spec)
	if err != nil {
		return err
	}

	results, err := testPaths(url, paths)
	if err != nil {
		return err
	}

	displayResultsAsTable(results)
	return nil
}

func loadAndValidateSpec(filePath string) (*openapi3.T, error) {
	loader := openapi3.NewLoader()
	spec, err := loader.LoadFromFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("loading OpenAPI spec: %w", err)
	}

	if err := spec.Validate(loader.Context); err != nil {
		return nil, fmt.Errorf("validating OpenAPI spec: %w", err)
	}

	return spec, nil
}

func printSpecInfo(spec *openapi3.T) {
	fmt.Println("Title:", spec.Info.Title)
	fmt.Println("Version:", spec.Info.Version)
}

func getServerURL(spec *openapi3.T) (string, error) {
	if len(spec.Servers) == 0 {
		return "", fmt.Errorf("no servers found in the specification")
	}
	return spec.Servers[0].URL, nil
}

func getPathsFromSpec(spec *openapi3.T) (map[string]interface{}, error) {
	data, err := json.MarshalIndent(spec.Paths, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshalling paths: %w", err)
	}

	var paths map[string]interface{}
	if err := json.Unmarshal(data, &paths); err != nil {
		return nil, fmt.Errorf("unmarshalling paths: %w", err)
	}

	return paths, nil
}

func testPaths(url string, paths map[string]interface{}) ([]TestResult, error) {
	var results []TestResult

	for endpoint, val := range paths {
		endpoint = strings.Replace(endpoint, "{resource}", "test", -1)
		endpoint = strings.Replace(endpoint, "{id}", "1", -1)

		methods, ok := val.(map[string]interface{})
		if !ok {
			continue
		}

		for method := range methods {
			result, err := testEndpoint(url, endpoint, method)
			if err != nil {
				return nil, err
			}
			results = append(results, result)
		}
	}

	sort.Slice(results, func(i, j int) bool {
		if results[i].Endpoint == results[j].Endpoint {
			return results[i].Method < results[j].Method
		}
		return results[i].Endpoint < results[j].Endpoint
	})

	return results, nil
}

func testEndpoint(url, endpoint, method string) (TestResult, error) {
	fullURL := url + endpoint
	ctx := context.Background()

	resp, err := makeRequest(ctx, strings.ToUpper(method), fullURL)
	status := "PASSED"

	if err != nil || resp.StatusCode >= 400 {
		status = "FAILED"
	}

	if resp != nil {
		_, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			status = "FAILED"
		}
	}

	return TestResult{Endpoint: endpoint, Method: strings.ToUpper(method), Status: status}, nil
}
