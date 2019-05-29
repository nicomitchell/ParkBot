package mocks

import (
	"fmt"
	"net/http"
)

//MockResponseWriter returns a mock http.ResponseWriter
type MockResponseWriter struct {
	StatusHeader int
	Message      []byte
	errFlag      bool
}

//Header will always return nil as it is not used
func (m *MockResponseWriter) Header() http.Header {
	return nil
}

func (m *MockResponseWriter) Write(b []byte) (int, error) {
	if m.errFlag {
		return 0, fmt.Errorf("test")
	}
	m.Message = append(m.Message, b...)
	return len(b), nil
}

//WriteHeader writes the header
func (m *MockResponseWriter) WriteHeader(statusCode int) {
	m.StatusHeader = statusCode
}

//NewMockResponseWriter returns a mock response writer
func NewMockResponseWriter(errFlag bool) http.ResponseWriter {
	return &MockResponseWriter{
		errFlag:      errFlag,
		Message:      []byte{},
		StatusHeader: 0,
	}
}
