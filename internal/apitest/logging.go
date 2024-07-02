// internal/apitest/logging.go

package apitest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
	"time"
)

type NetworkLog struct {
	Timestamp time.Time   `json:"timestamp"`
	Request   RequestLog  `json:"request"`
	Response  ResponseLog `json:"response"`
}

type RequestLog struct {
	Method  string      `json:"method"`
	URL     string      `json:"url"`
	Headers http.Header `json:"headers"`
	Body    string      `json:"body,omitempty"`
}

type ResponseLog struct {
	StatusCode int         `json:"status_code"`
	Status     string      `json:"status"`
	Headers    http.Header `json:"headers"`
	Body       string      `json:"body"`
}

var (
	logMutex sync.Mutex
	logFile  *os.File
	logs     []NetworkLog
)

func InitializeLogFile(logFilePath string) error {
	var err error
	logFile, err = os.Create(logFilePath) // This will create a new file or truncate an existing one
	if err != nil {
		return fmt.Errorf("creating log file: %w", err)
	}

	logs = []NetworkLog{} // Initialize an empty slice for logs
	return nil
}

func CloseLogFile() error {
	if logFile != nil {
		logMutex.Lock()
		defer logMutex.Unlock()

		// Marshal all logs
		content, err := json.MarshalIndent(logs, "", "  ")
		if err != nil {
			return fmt.Errorf("marshaling logs: %w", err)
		}

		// Write all logs to file
		_, err = logFile.Write(content)
		if err != nil {
			return fmt.Errorf("writing to log file: %w", err)
		}

		return logFile.Close()
	}
	return nil
}

func LogNetworkInteraction(req *http.Request, resp *http.Response) error {
	networkLog := NetworkLog{
		Timestamp: time.Now(),
		Request: RequestLog{
			Method:  req.Method,
			URL:     req.URL.String(),
			Headers: req.Header,
		},
		Response: ResponseLog{
			StatusCode: resp.StatusCode,
			Status:     resp.Status,
			Headers:    resp.Header,
		},
	}

	// Log request body if present
	if req.Body != nil {
		bodyBytes, err := io.ReadAll(req.Body)
		if err != nil {
			return fmt.Errorf("reading request body: %w", err)
		}
		req.Body.Close()
		req.Body = io.NopCloser(bytes.NewReader(bodyBytes))
		networkLog.Request.Body = string(bodyBytes)
	}

	// Log response body
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("reading response body: %w", err)
	}
	resp.Body.Close()
	resp.Body = io.NopCloser(bytes.NewReader(bodyBytes))
	networkLog.Response.Body = string(bodyBytes)

	logMutex.Lock()
	logs = append(logs, networkLog)
	logMutex.Unlock()

	return nil
}
