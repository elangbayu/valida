package apitest

import "github.com/getkin/kin-openapi/openapi3"

type Server struct {
	URL string `json:"url"`
}

type TestResult struct {
	Endpoint string
	Method   string
	Status   string
}

type EndpointInfo struct {
	Method      string
	Parameters  []*openapi3.Parameter
	RequestBody *openapi3.RequestBody
}
