package apitest

import (
	"fmt"

	"github.com/getkin/kin-openapi/openapi3"
)

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
	fmt.Printf("Title: %s\n", spec.Info.Title)
	fmt.Printf("Version: %s\n", spec.Info.Version)
}

func getBaseURL(spec *openapi3.T) (string, error) {
	if len(spec.Servers) == 0 {
		return "", fmt.Errorf("no servers found in the specification")
	}
	return spec.Servers[0].URL, nil
}
