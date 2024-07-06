package apitest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/getkin/kin-openapi/openapi3"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"regexp"
	"strconv"
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
				// jsonBody := []byte(`{"email": "elang@mail.com", "password": "password"}`)
				// Dynamic
				reqBody := make(map[string]interface{})

				var mapss openapi3.Schemas
				switch strings.ToUpper(operation.Method) {
				case "GET":
					mapss = apiSpec.Spec.Paths.Find(pathItem.Path).Get.RequestBody.Value.Content.Get("application/json").Schema.Value.Properties
				case "POST":
					mapss = apiSpec.Spec.Paths.Find(pathItem.Path).Post.RequestBody.Value.Content.Get("application/json").Schema.Value.Properties
				case "PUT":
					mapss = apiSpec.Spec.Paths.Find(pathItem.Path).Put.RequestBody.Value.Content.Get("application/json").Schema.Value.Properties
				case "DELETE":
					mapss = apiSpec.Spec.Paths.Find(pathItem.Path).Delete.RequestBody.Value.Content.Get("application/json").Schema.Value.Properties
				case "PATCH":
					mapss = apiSpec.Spec.Paths.Find(pathItem.Path).Patch.RequestBody.Value.Content.Get("application/json").Schema.Value.Properties
				}

				for k, v := range mapss {
					pattern := `\b(?:e[-]?mail|mail)\b`
					re := regexp.MustCompile(pattern)
					matches := re.FindAllString(k, -1)
					switch v.Value.Type.Slice()[0] {
					case "string":
						if len(matches) > 0 {
							reqBody[k] = FakeEmail()
						} else {
							reqBody[k] = FakeString()
						}
					case "int":
						reqBody[k] = FakeInt()
					}
				}

				b, _ := json.Marshal(reqBody)

				jsonReader = bytes.NewReader(b)
			}

			req, err := http.NewRequest(strings.ToUpper(operation.Method), apiSpec.BaseURL+pathItem.Path, jsonReader)
			if err != nil {
				log.Fatalf("Error making request: %v", err)
			}
			req.Header.Add("Content-Type", "application/json")
			fmt.Println("Req ", apiSpec.BaseURL+pathItem.Path)

			if operation.Parameters != nil {
				fmt.Println("Operations Parameters: ", operation.Parameters)
				q := req.URL.Query()
				for _, param := range operation.Parameters {
					fmt.Println("Param: ", param)
					// Get Value "in"
					if inValue, ok := param["in"]; ok {
						fmt.Println("In: ", inValue)
						schema := param["schema"].(map[string]interface{})
						typeValue, _ := schema["type"]

						getValue := func(typeValue interface{}) string {
							if typeValue == "string" {
								return FakeString()
							} else if typeValue == "integer" {
								return strconv.Itoa(FakeInt())
							}
							return ""
						}

						value := getValue(typeValue)
						paramName := param["name"].(string)

						switch inValue {
						case "query":
							fmt.Println("Value Param: ", value)
							q.Add(paramName, value)
						case "header":
							fmt.Println("Value Header: ", value)
							req.Header.Add(paramName, value)
						}
					}
				}
				req.URL.RawQuery = q.Encode()
			}

			// Print Request
			requestDump, err := httputil.DumpRequest(req, true)
			if err != nil {
				log.Fatalf("Error dumping request: %v", err)
			}
			fmt.Println("MAKE Request: ", string(requestDump))
			fmt.Println("-------------------------")

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

// replaceDynamicPlaceholders replaces placeholders enclosed in curly braces with random strings.
func replaceDynamicPlaceholders(path string) string {
	re := regexp.MustCompile(`{[^{}]*}`) // Regular expression to find words within curly braces
	return re.ReplaceAllStringFunc(path, func(match string) string {
		// Generate a random string of length 10 for each match (this length can be adjusted)
		return FakeString()
	})
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
