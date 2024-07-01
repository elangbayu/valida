package apitest

type Server struct {
	URL string `json:"url"`
}

type TestResult struct {
	Endpoint string
	Method   string
	Status   string
}
