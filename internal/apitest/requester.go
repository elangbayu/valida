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
var logger *Logger

func MakeRequest(apiSpec *APISpec) {
	var err error
	logger, err = NewLogger()
	if err != nil {
		fmt.Printf("Failed to create logger: %v\n", err)
		return
	}
	defer logger.Close()

	var tableRows []TableRow

	for _, pathItem := range apiSpec.Paths {
		for _, operation := range pathItem.Operations {
			endpoint := apiSpec.BaseURL + pathItem.Path
			method := strings.ToUpper(operation.Method)

			req, requestBody := prepareRequest(apiSpec, pathItem, operation)
			if req == nil {
				logger.LogError(fmt.Errorf("failed to prepare request for %s %s", method, endpoint))
				tableRows = append(tableRows, TableRow{
					Endpoint:  endpoint,
					Method:    method,
					Response:  "N/A",
					Assertion: "FAIL: Request preparation error",
				})
				continue
			}

			logger.LogRequest(req, requestBody)

			expectedResponse, _ := GetExpectedResponse(operation)
			resp, responseBody, assertionResult := requestAndValidate(req, expectedResponse, endpoint, method)

			if resp != nil {
				logger.LogResponse(resp, responseBody)
				responseInfo := fmt.Sprintf("%d %s", resp.StatusCode, http.StatusText(resp.StatusCode))
				tableRows = append(tableRows, TableRow{
					Endpoint:  endpoint,
					Method:    method,
					Response:  responseInfo,
					Assertion: assertionResult,
				})
			} else {
				logger.LogError(fmt.Errorf("no response received for %s %s", method, endpoint))
				tableRows = append(tableRows, TableRow{
					Endpoint:  endpoint,
					Method:    method,
					Response:  "No response",
					Assertion: assertionResult,
				})
			}
		}
	}

	DisplayTable(tableRows)
}

func prepareRequest(apiSpec *APISpec, pathItem *PathItem, operation *Operation) (*http.Request, string) {
	var jsonReader io.Reader
	var requestBody string

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
		requestBody = string(b)
	}

	// Replace path parameters with fake values
	endpoint := replacePathParameters(apiSpec.BaseURL+pathItem.Path, operation.Parameters)

	req, err := http.NewRequest(strings.ToUpper(operation.Method), endpoint, jsonReader)
	if err != nil {
		return nil, ""
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

	return req, requestBody
}

func replacePathParameters(path string, parameters []map[string]interface{}) string {
	for _, param := range parameters {
		if in, ok := param["in"].(string); ok && in == "path" {
			name, ok := param["name"].(string)
			if !ok {
				continue
			}
			placeholder := fmt.Sprintf("{%s}", name)
			var fakeValue string
			if name == "id" {
				fakeValue = strconv.Itoa(randInt(1, 10))
			} else {
				fakeValue = generateFakeValue(param["schema"].(map[string]interface{}))
			}
			path = strings.Replace(path, placeholder, fakeValue, 1)
		}
	}
	return path
}

func generateFakeValue(schema map[string]interface{}) string {
	schemaType, ok := schema["type"].(string)
	if !ok {
		return FakeString()
	}

	switch schemaType {
	case "string":
		return FakeString()
	case "integer":
		return strconv.Itoa(1)
	case "number":
		return fmt.Sprintf("%.2f", float64(FakeInt()))
	default:
		return FakeString()
	}
}

func requestAndValidate(req *http.Request, expectedResp *ExpectedResponse, endpoint, method string) (*http.Response, string, string) {
	resp, err := client.Do(req)
	if err != nil {
		logger.LogError(fmt.Errorf("error doing request: %v", err))
		return nil, "", fmt.Sprintf("FAIL: Error doing request: %v", err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.LogError(fmt.Errorf("error reading body: %v", err))
		return resp, "", fmt.Sprintf("FAIL: Error reading body: %v", err)
	}

	responseBody := string(body)

	if expectedResp != nil {
		if err := CompareResponses(resp, expectedResp); err != nil {
			return resp, responseBody, fmt.Sprintf("FAIL: %v", err)
		} else {
			return resp, responseBody, "PASS"
		}
	} else {
		return resp, responseBody, "WARNING: No expected response to validate against"
	}
}
