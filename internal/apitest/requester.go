package apitest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

var client = &http.Client{}

func MakeRequest(apiSpec *APISpec) {
	fmt.Printf("Base URL: %s\n", apiSpec.BaseURL)

	for _, pathItem := range apiSpec.Paths {
		for _, operation := range pathItem.Operations {
			PrintEndpointInfo(apiSpec.BaseURL, pathItem.Path, operation.Method)

			var jsonReader io.Reader
			if operation.RequestBody != nil {
				reqBody := make(map[string]interface{})

				properties := operation.RequestBody["content"].(map[string]interface{})["application/json"].(map[string]interface{})["schema"].(map[string]interface{})["properties"].(map[string]interface{})

				for k, v := range properties {
					pattern := `\b(?:e[-]?mail|mail)\b`
					re := regexp.MustCompile(pattern)
					matches := re.FindAllString(k, -1)
					vType := v.(map[string]interface{})["type"].(string)
					switch vType {
					case "string":
						if len(matches) > 0 {
							reqBody[k] = FakeEmail()
						} else {
							reqBody[k] = FakeString()
						}
					case "integer":
						reqBody[k] = FakeInt()
					}
				}

				b, _ := json.Marshal(reqBody)
				jsonReader = bytes.NewReader(b)
				PrintRequestBody(b)
			}

			req, err := http.NewRequest(strings.ToUpper(operation.Method), apiSpec.BaseURL+pathItem.Path, jsonReader)
			if err != nil {
				PrintError(apiSpec.BaseURL+pathItem.Path, operation.Method, fmt.Sprintf("Error making request: %v", err))
				continue
			}
			req.Header.Add("Content-Type", "application/json")

			if operation.Parameters != nil {
				q := req.URL.Query()
				for _, param := range operation.Parameters {
					if inValue, ok := param["in"]; ok {
						schema := param["schema"].(map[string]interface{})
						typeValue, _ := schema["type"]

						getValue := func(paramName string, typeValue interface{}) string {
							if paramName == "page" {
								return "1"
							}
							if typeValue == "string" {
								return FakeString()
							} else if typeValue == "integer" {
								return strconv.Itoa(FakeInt())
							}
							return ""
						}

						value := getValue(param["name"].(string), typeValue)
						paramName := param["name"].(string)

						switch inValue {
						case "query":
							q.Add(paramName, value)
						case "header":
							req.Header.Add(paramName, value)
						}
					}
				}
				req.URL.RawQuery = q.Encode()
			}

			PrintRequest(req)

			// Do Req and Validate Response
			expectedResponse, err := GetExpectedResponse(operation)
			if err != nil {
				PrintError(apiSpec.BaseURL+pathItem.Path, operation.Method, fmt.Sprintf("Error getting expected response: %v", err))
				// Set expectedResponse to nil to indicate that no validation will be done
				expectedResponse = nil
			}

			requestAndValidate(req, expectedResponse, apiSpec.BaseURL+pathItem.Path, operation.Method)
		}
	}
}

func requestAndValidate(req *http.Request, expectedResp *ExpectedResponse, endpoint, method string) {
	resp, err := client.Do(req)
	if err != nil {
		PrintError(endpoint, method, fmt.Sprintf("Error doing request: %v", err))
		return
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			PrintError(endpoint, method, fmt.Sprintf("Error closing response body: %v", err))
		}
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		PrintError(endpoint, method, fmt.Sprintf("Error reading body: %v", err))
		return
	}

	PrintResponse(resp, body)

	if expectedResp != nil {
		if err := CompareResponses(resp, expectedResp); err != nil {
			PrintAssertionResult(err)
		} else {
			PrintAssertionResult(nil)
		}
	} else {
		fmt.Println("Assertion Result: ⚠️ No expected response to validate against.")
	}
}
