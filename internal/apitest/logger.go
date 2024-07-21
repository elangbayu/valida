package apitest

import (
	"fmt"
	"net/http"
	"os"
	"time"
)

type Logger struct {
	file *os.File
}

func NewLogger() (*Logger, error) {
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	filename := fmt.Sprintf("test_log_%s.txt", timestamp)

	file, err := os.Create(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to create log file: %w", err)
	}

	return &Logger{file: file}, nil
}

func (l *Logger) Close() error {
	return l.file.Close()
}

func (l *Logger) LogRequest(req *http.Request, body string) {
	l.writeLog(fmt.Sprintf("REQUEST: %s %s\n", req.Method, req.URL))
	l.writeLog("Headers:\n")
	for key, values := range req.Header {
		for _, value := range values {
			l.writeLog(fmt.Sprintf("%s: %s\n", key, value))
		}
	}
	if body != "" {
		l.writeLog(fmt.Sprintf("Body:\n%s\n", body))
	}
	l.writeLog("\n")
}

func (l *Logger) LogResponse(resp *http.Response, body string) {
	l.writeLog(fmt.Sprintf("RESPONSE: %s\n", resp.Status))
	l.writeLog("Headers:\n")
	for key, values := range resp.Header {
		for _, value := range values {
			l.writeLog(fmt.Sprintf("%s: %s\n", key, value))
		}
	}
	l.writeLog(fmt.Sprintf("Body:\n%s\n", body))
	l.writeLog("\n")
}

func (l *Logger) LogError(err error) {
	l.writeLog(fmt.Sprintf("ERROR: %v\n\n", err))
}

func (l *Logger) writeLog(message string) {
	_, err := l.file.WriteString(message)
	if err != nil {
		fmt.Printf("Failed to write to log file: %v\n", err)
	}
}
