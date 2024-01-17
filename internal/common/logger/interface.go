package logger

import "net/http"

// BaseLogger used logger interface definition.
type BaseLogger interface {
	// Sync Method for flushing data in stream.
	Sync() error
	// Infof posts message on log 'info' Level.
	Infof(msg string, fields ...interface{})
	// InfoWithFieldsf posts message on log 'info' Level with extra fields.
	InfoWithFieldsf(fields map[string]interface{}, msg string, msgFields ...interface{})
	// Fatalf posts message on log 'fatal' Level.
	Fatalf(msg string, fields ...interface{})
	// LogHandler implements middleware for logging requests.
	LogHandler(h http.Handler) http.Handler
}
