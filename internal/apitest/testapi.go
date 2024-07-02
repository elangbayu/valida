package apitest

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

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

func getPathsFromSpec(spec *openapi3.T) (map[string]map[string]EndpointInfo, error) {
	paths := make(map[string]map[string]EndpointInfo)

	for path, pathItem := range spec.Paths.Map() {
		endpoints := make(map[string]EndpointInfo)

		operations := map[string]*openapi3.Operation{
			"GET":    pathItem.Get,
			"POST":   pathItem.Post,
			"PUT":    pathItem.Put,
			"DELETE": pathItem.Delete,
			"PATCH":  pathItem.Patch,
		}

		for method, operation := range operations {
			if operation != nil {
				var requestBody *openapi3.RequestBody
				if operation.RequestBody != nil {
					requestBody = operation.RequestBody.Value
				}
				parameters := dereferenceParams(append(pathItem.Parameters, operation.Parameters...))
				endpointInfo := EndpointInfo{
					Method:      method,
					Parameters:  parameters,
					RequestBody: requestBody,
				}
				endpoints[method] = endpointInfo
			}
		}

		paths[path] = endpoints
	}

	return paths, nil
}

func dereferenceParams(params openapi3.Parameters) []*openapi3.Parameter {
	result := make([]*openapi3.Parameter, 0, len(params))
	for _, param := range params {
		if param != nil {
			result = append(result, param.Value)
		}
	}
	return result
}

func testPaths(url string, paths map[string]map[string]EndpointInfo) ([]TestResult, error) {
	var results []TestResult

	for path, methods := range paths {
		for method, info := range methods {
			result, err := testEndpoint(url, path, method, info)
			if err != nil {
				return nil, err
			}
			results = append(results, result)
		}
	}

	return results, nil
}

func testEndpoint(baseURL, path, method string, info EndpointInfo) (TestResult, error) {
	fullURL := baseURL + path
	req, err := constructRequest(method, fullURL, info)
	if err != nil {
			return TestResult{}, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := makeRequest(ctx, req.Method, req.URL.String())
	if err != nil {
			return TestResult{}, err
	}
	defer resp.Body.Close()

	// Log the network interaction
	err = LogNetworkInteraction(req, resp)
	if err != nil {
			fmt.Printf("Error logging network interaction: %v\n", err)
	}

	return TestResult{
			Endpoint: path,
			Method:   method,
			Status:   resp.Status,
	}, nil
}

func constructRequest(method, url string, info EndpointInfo) (*http.Request, error) {
	var body io.Reader
	if info.RequestBody != nil {
		// Here you would generate a valid request body based on the schema
		// For now, we'll just use an empty body
		body = strings.NewReader("{}")
	}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	// Add query parameters
	q := req.URL.Query()
	for _, param := range info.Parameters {
		if param.In == "query" {
			// Here you would generate a valid value based on the parameter schema
			// For now, we'll just use a placeholder value
			q.Add(param.Name, "test_value")
		}
	}
	req.URL.RawQuery = q.Encode()

	// Add headers
	for _, param := range info.Parameters {
		if param.In == "header" {
			// Here you would generate a valid value based on the parameter schema
			// For now, we'll just use a placeholder value
			req.Header.Add(param.Name, "test_value")
		}
	}

	return req, nil
}
