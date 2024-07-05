package apitest

import (
	"fmt"
	"strings"
)

func PrintAPISpec(apiSpec *APISpec) {
	fmt.Printf("Base URL: %s\n\n", apiSpec.BaseURL)

	for _, pathItem := range apiSpec.Paths {
		fmt.Printf("Path: %s\n", pathItem.Path)
		for _, operation := range pathItem.Operations {
			fmt.Printf("  Method: %s\n", operation.Method)
			printParameters(operation.Parameters)
			printRequestBody(operation.RequestBody)
			printResponses(operation.Responses)
			fmt.Println()
		}
		fmt.Println()
	}
}

func printParameters(parameters []map[string]interface{}) {
	if len(parameters) > 0 {
		fmt.Println("  Parameters:")
		for i, param := range parameters {
			fmt.Printf("    Parameter %d:\n", i+1)
			printNestedMap(param, 3)
		}
	}
}

func printRequestBody(requestBody map[string]interface{}) {
	if len(requestBody) > 0 {
		fmt.Println("  Request Body:")
		printNestedMap(requestBody, 2)
	}
}

func printResponses(responses map[string]map[string]interface{}) {
	if len(responses) > 0 {
		fmt.Println("  Responses:")
		for statusCode, response := range responses {
			fmt.Printf("    Status Code %s:\n", statusCode)
			printNestedMap(response, 3)
		}
	}
}

func printNestedMap(m map[string]interface{}, indent int) {
	for key, value := range m {
		fmt.Printf("%s%s: ", strings.Repeat("  ", indent), key)
		switch v := value.(type) {
		case map[string]interface{}:
			fmt.Println()
			printNestedMap(v, indent+1)
		case []interface{}:
			fmt.Println()
			for _, item := range v {
				if itemMap, ok := item.(map[string]interface{}); ok {
					printNestedMap(itemMap, indent+1)
				} else {
					fmt.Printf("%s- %v\n", strings.Repeat("  ", indent+1), item)
				}
			}
		default:
			fmt.Printf("%v\n", v)
		}
	}
}
