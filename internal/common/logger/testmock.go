package logger

import (
	"net/http"
)

type logMock struct {
}

// CreateMock creates debug plug for ignoring logger in test.
func CreateMock() BaseLogger {
	return &logMock{}
}

func (t *logMock) Infof(_ string, _ ...interface{}) {
}

func (t *logMock) InfoWithFieldsf(_ map[string]interface{}, _ string, _ ...interface{}) {
}

func (t *logMock) Printf(_ string, _ ...interface{}) {
}

func (t *logMock) Fatal(_ ...interface{}) {
}

func (t *logMock) Fatalf(_ string, _ ...interface{}) {
}

func (t *logMock) Sync() error {
	return nil
}

func (t *logMock) LogHandler(h http.Handler) http.Handler {
	return h
}
