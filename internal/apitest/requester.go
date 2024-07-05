package apitest

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

var client = &http.Client{}

func MakeRequest(apiSpec *APISpec) {
	fmt.Printf("Base URL: %s\n\n", apiSpec.BaseURL)

	for _, pathItem := range apiSpec.Paths {
		fmt.Printf("Path: %s\n", pathItem.Path)
		for _, operation := range pathItem.Operations {
			fmt.Printf("  Method: %s\n", operation.Method)
			fmt.Printf("  Req Body: %s\n", operation.RequestBody)

			var jsonReader io.Reader
			if operation.RequestBody != nil {
				// tambahkan request body ke request
				jsonBody := []byte(`{"email": "elang@mail.com", "password": "password"}`)
				jsonReader = bytes.NewReader(jsonBody)
			}

			// http request ke URL: apiSpec.BaseURL + pathItem.Path, methodnya pake operation.Method
			req, err := http.NewRequest(strings.ToUpper(operation.Method), apiSpec.BaseURL+pathItem.Path, jsonReader)
			if err != nil {
				log.Fatalf("Error making request: %v", err)
			}
			req.Header.Add("Content-Type", "application/json")
			fmt.Println("Req Method", operation.Method)
			fmt.Println("Req ", apiSpec.BaseURL+pathItem.Path)

			if operation.Parameters != nil {
				fmt.Println("Operations Parameters: ", operation.Parameters)
				q := req.URL.Query()
				for key, param := range operation.Parameters {
					fmt.Println("Key Param: ", key, param)
					// if param.In == "query" {
					// 	// Here you would generate a valid value based on the parameter schema
					// 	// For now, we'll just use a placeholder value
					// 	q.Add(param.Name, "test_value")
					// }
				}
				req.URL.RawQuery = q.Encode()
			}

			// Mencetak request sebelum mengirimkannya
			// requestDump, err := httputil.DumpRequest(req, true)
			// if err != nil {
			// 	log.Fatalf("Error dumping request: %v", err)
			// }
			// fmt.Println("Request: ", string(requestDump))

			// Do Req
			request(req)

			fmt.Println()
		}
		fmt.Println()
	}
}

func request(req *http.Request) {
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error doing request: %v", err)
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Fatalf("Error closing response body: %v", err)
		}
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading body: %v", err)
	}

	fmt.Println("Response Status:", resp.Status)
	fmt.Println("Response Body:", string(body))
}

// func printParameters(parameters []map[string]interface{}) {
// 	if len(parameters) > 0 {
// 		fmt.Println("  Parameters:")
// 		for i, param := range parameters {
// 			fmt.Printf("    Parameter %d:\n", i+1)
// 			printNestedMap(param, 3)
// 		}
// 	}
// }
//
// func printRequestBody(requestBody map[string]interface{}) {
// 	if len(requestBody) > 0 {
// 		fmt.Println("  Request Body:")
// 		printNestedMap(requestBody, 2)
// 	}
// }
//
// func printResponses(responses map[string]map[string]interface{}) {
// 	if len(responses) > 0 {
// 		fmt.Println("  Responses:")
// 		for statusCode, response := range responses {
// 			fmt.Printf("    Status Code %s:\n", statusCode)
// 			printNestedMap(response, 3)
// 		}
// 	}
// }
//
// func printNestedMap(m map[string]interface{}, indent int) {
// 	for key, value := range m {
// 		fmt.Printf("%s%s: ", strings.Repeat("  ", indent), key)
// 		switch v := value.(type) {
// 		case map[string]interface{}:
// 			fmt.Println()
// 			printNestedMap(v, indent+1)
// 		case []interface{}:
// 			fmt.Println()
// 			for _, item := range v {
// 				if itemMap, ok := item.(map[string]interface{}); ok {
// 					printNestedMap(itemMap, indent+1)
// 				} else {
// 					fmt.Printf("%s- %v\n", strings.Repeat("  ", indent+1), item)
// 				}
// 			}
// 		default:
// 			fmt.Printf("%v\n", v)
// 		}
// 	}
// }
