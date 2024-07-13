package apitest

import (
	"fmt"

	"github.com/getkin/kin-openapi/openapi3"
)

// APISpec represents the parsed API specification
type APISpec struct {
	Spec    *openapi3.T
	BaseURL string
	Paths   map[string]*PathItem
}

// TestAPISpec is the main function to test the API specification
func TestAPISpec(filePath string) (*APISpec, error) {
	spec, err := loadAndValidateSpec(filePath)
	if err != nil {
		return nil, fmt.Errorf("error validating OpenAPI spec: %w", err)
	}

	if err := loadConfig(filePath); err != nil {
		return nil, fmt.Errorf("error reading OpenAPI spec: %w", err)
	}

	apiSpec := &APISpec{
		Spec:  spec,
		Paths: make(map[string]*PathItem),
	}

	printSpecInfo(apiSpec.Spec)

	baseURL, err := getBaseURL(apiSpec.Spec)
	if err != nil {
		return nil, fmt.Errorf("baseURL not found: %w", err)
	}
	apiSpec.BaseURL = baseURL

	if err := processPaths(apiSpec); err != nil {
		return nil, fmt.Errorf("error processing paths: %w", err)
	}

	return apiSpec, nil
}
