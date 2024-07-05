package apitest

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/spf13/viper"
)

// PathItem represents a path in the API specification
type PathItem struct {
	Path       string
	Operations map[string]*Operation
}

// Operation represents an operation in the API specification
type Operation struct {
	Method      string
	Parameters  []map[string]interface{}
	RequestBody map[string]interface{}
	Responses   map[string]map[string]interface{}
}

func processPaths(apiSpec *APISpec) error {
	paths := viper.GetStringMap("paths")

	for path, pathMap := range paths {
		pathItem := &PathItem{
			Path:       replaceDynamicPlaceholders(path),
			Operations: make(map[string]*Operation),
		}

		// path -> operation -> { parameters, request body dan responses }

		pathMaps, ok := pathMap.(map[string]interface{})
		if !ok {
			return fmt.Errorf("invalid path map for path: %s", path)
		}

		if err := processOperations(pathItem, pathMaps); err != nil {
			return err
		}

		apiSpec.Paths[path] = pathItem
	}

	return nil
}

// replaceDynamicPlaceholders replaces placeholders enclosed in curly braces with random strings.
func replaceDynamicPlaceholders(path string) string {
	re := regexp.MustCompile(`{[^{}]*}`) // Regular expression to find words within curly braces
	return re.ReplaceAllStringFunc(path, func(match string) string {
		// Generate a random string of length 10 for each match (this length can be adjusted)
		return FakeString()
	})
}

func processOperations(pathItem *PathItem, pathMaps map[string]interface{}) error {
	for method, operationMap := range pathMaps {
		operation := &Operation{
			Method: method,
		}

		operationMaps, ok := operationMap.(map[string]interface{})
		if !ok {
			return fmt.Errorf("invalid operation map for operation: %s", method)
		}

		if err := processOperationDetails(operation, operationMaps); err != nil {
			return err
		}

		pathItem.Operations[method] = operation
	}

	return nil
}

func processOperationDetails(operation *Operation, operationMaps map[string]interface{}) error {
	for k, v := range operationMaps {
		switch strings.ToLower(k) {
		case "parameters":
			if params, ok := v.([]interface{}); ok {
				operation.Parameters = make([]map[string]interface{}, len(params))
				for i, param := range params {
					if paramMap, ok := param.(map[string]interface{}); ok {
						operation.Parameters[i] = paramMap
					}
				}
			}
		case "requestbody":
			if reqBody, ok := v.(map[string]interface{}); ok {
				operation.RequestBody = reqBody
			}
		case "responses":
			if responses, ok := v.(map[string]interface{}); ok {
				operation.Responses = make(map[string]map[string]interface{})
				for statusCode, response := range responses {
					if respMap, ok := response.(map[string]interface{}); ok {
						operation.Responses[statusCode] = respMap
					}
				}
			}
		}
	}

	return nil
}
