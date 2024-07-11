package apitest

import (
	"fmt"
	"net/http"
	"net/http/httputil"
)

// PrintEndpointInfo prints the endpoint and method information
func PrintEndpointInfo(baseURL, path, method string) {
	fmt.Println("========================================")
	fmt.Printf("Endpoint: %s\n", baseURL+path)
	fmt.Printf("Method: %s\n", method)
	fmt.Println("----------------------------------------")
}

// PrintRequestBody prints the request body
func PrintRequestBody(body []byte) {
	fmt.Println("Request Body:")
	fmt.Printf("%s\n", body)
	fmt.Println("----------------------------------------")
}

// PrintRequest prints the request details
func PrintRequest(req *http.Request) {
	requestDump, err := httputil.DumpRequest(req, true)
	if err != nil {
		fmt.Printf("Error dumping request: %v\n", err)
		return
	}
	fmt.Println("Request:")
	fmt.Printf("%s\n", string(requestDump))
	fmt.Println("----------------------------------------")
}

// PrintResponse prints the response details
func PrintResponse(resp *http.Response, body []byte) {
	fmt.Println("Response:")
	fmt.Printf("Status: %d\n", resp.StatusCode)
	fmt.Printf("Body: %s\n", string(body))
	fmt.Println("----------------------------------------")
}

// PrintAssertionResult prints the result of the assertion
func PrintAssertionResult(err error) {
	if err != nil {
		fmt.Printf("Assertion Result: ❌ Response validation failed: %v\n", err)
	} else {
		fmt.Println("Assertion Result: ✔️ Response validation passed!")
	}
	fmt.Println("========================================\n")
}

// PrintError prints the error for a specific endpoint and method
func PrintError(endpoint, method, errorMessage string) {
	fmt.Println("========================================")
	fmt.Printf("Endpoint: %s\n", endpoint)
	fmt.Printf("Method: %s\n", method)
	fmt.Println("Error:")
	fmt.Printf("%s\n", errorMessage)
	fmt.Println("========================================\n")
}
