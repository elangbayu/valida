package apitest

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
)

// ExpectedResponse represents the expected response based on OpenAPI specification
type ExpectedResponse struct {
	StatusCode int
	Body       map[string]interface{}
}

// GetExpectedResponse extracts the expected response for a given operation from the OpenAPI specification
func GetExpectedResponse(operation *Operation) (*ExpectedResponse, error) {
	response, ok := operation.Responses["200"]
	if !ok {
		return nil, fmt.Errorf("no 200 response found in the specification")
	}

	content, ok := response["content"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("no content found in the 200 response")
	}

	mediaType, ok := content["application/json"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("no application/json media type found in the 200 response")
	}

	schema, ok := mediaType["schema"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("no schema found in the application/json media type")
	}

	example, ok := schema["example"]
	if !ok {
		return nil, fmt.Errorf("no example found in the schema")
	}

	expectedBody := make(map[string]interface{})
	if err := json.Unmarshal([]byte(example.(string)), &expectedBody); err != nil {
		return nil, fmt.Errorf("error unmarshalling expected body: %w", err)
	}

	return &ExpectedResponse{
		StatusCode: 200,
		Body:       expectedBody,
	}, nil
}

// CompareResponses compares the actual response with the expected response
func CompareResponses(actualResp *http.Response, expectedResp *ExpectedResponse) error {
	if actualResp.StatusCode != expectedResp.StatusCode {
		return fmt.Errorf("expected status code %d, but got %d", expectedResp.StatusCode, actualResp.StatusCode)
	}

	actualBody, err := io.ReadAll(actualResp.Body)
	if err != nil {
		return fmt.Errorf("error reading response body: %w", err)
	}
	defer actualResp.Body.Close()

	var actualBodyJSON map[string]interface{}
	if err := json.Unmarshal(actualBody, &actualBodyJSON); err != nil {
		return fmt.Errorf("error unmarshalling actual response body: %w", err)
	}

	if !reflect.DeepEqual(expectedResp.Body, actualBodyJSON) {
		return fmt.Errorf("expected response body %v, but got %v", expectedResp.Body, actualBodyJSON)
	}

	return nil
}
