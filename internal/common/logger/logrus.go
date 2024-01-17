package logger

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

var (
	_ BaseLogger = (*Logrus)(nil)
)

// Config is a set of params used in constructor.
type Config struct {
	Level string
	File  string
}

// Logrus BaseLogger implementation based on logrus.
type Logrus struct {
	file *os.File
}

// NewLogrus returns def logger
func NewLogrus(cfg *Config) (BaseLogger, error) {
	level, err := logrus.ParseLevel(cfg.Level)
	if err != nil {
		return nil, fmt.Errorf("parse Level from config: %w", err)
	}
	logrus.SetLevel(level)

	logrus.SetFormatter(&logrus.JSONFormatter{})

	var file *os.File
	if cfg.File != "" {
		file, err = os.OpenFile(cfg.File, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			return nil, fmt.Errorf("open File '%s' for logging: %w", cfg.File, err)
		}

		logrus.SetOutput(file)
	} else {
		logrus.SetOutput(os.Stdout)
	}

	return &Logrus{
		file: file,
	}, nil
}

func (l *Logrus) Infof(msg string, fields ...interface{}) {
	logrus.Infof(msg, fields...)
}

func (l *Logrus) InfoWithFieldsf(fields map[string]interface{}, msg string, msgFields ...interface{}) {
	logrus.WithFields(fields).Infof(msg, msgFields...)
}

func (l *Logrus) Fatal(msg ...interface{}) {
	logrus.Fatal(msg...)
}

func (l *Logrus) Fatalf(msg string, fields ...interface{}) {
	logrus.Fatalf(msg, fields...)
}

func (l *Logrus) Sync() error {
	if l.file != nil {
		logrus.Fatal(l.file.Close())
	}

	return nil
}

func (l *Logrus) LogHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		loggingWriter := createResponseWriter(w)
		h.ServeHTTP(loggingWriter, r)
		duration := time.Since(start)

		logrus.WithFields(logrus.Fields{
			"uri":              r.RequestURI,
			"method":           r.Method,
			"status":           loggingWriter.getResponseData().status,
			"content-type":     loggingWriter.Header().Get("Content-Type"),
			"content-encoding": loggingWriter.Header().Get("Content-Encoding"),
			"HashSHA256":       loggingWriter.Header().Get("HashSHA256"),
			"duration":         duration.String(),
			"size":             loggingWriter.getResponseData().size,
		}).Info("new incoming HTTP request")
	})
}
