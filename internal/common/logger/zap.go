package logger

import (
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	_ BaseLogger = (*Zap)(nil)
)

// Zap wrapper of Zap logger.
type Zap struct {
	zap *zap.Logger
}

// NewZap create method for zap logger.
func NewZap(level string) (BaseLogger, error) {
	cfg, err := initConfig(level)
	if err != nil {
		return nil, err
	}

	logZap, err := cfg.Build()
	if err != nil {
		return nil, fmt.Errorf("create zap loggerZap^ %w", err)
	}

	return &Zap{zap: logZap}, nil
}

func (l *Zap) Sync() error {
	return l.zap.Sync()
}

// Infof generates 'info' level log.
func (l *Zap) Infof(msg string, fields ...interface{}) {
	l.zap.Info(fmt.Sprintf(msg, fields...))
}

// InfoWithFieldsf generates 'info' level log.
func (l *Zap) InfoWithFieldsf(fields map[string]interface{}, msg string, msgFields ...interface{}) {
	var zapFields []zap.Field
	for k, v := range fields {
		zapFields = append(zapFields, zap.Field{
			Key:    k,
			Type:   zapcore.StringType,
			String: fmt.Sprint(v),
		})
	}
	l.zap.With(zapFields...).Info(fmt.Sprintf(msg, msgFields...))
}

// Fatalf generates 'info' level log.
func (l *Zap) Fatalf(msg string, fields ...interface{}) {
	l.zap.Fatal(fmt.Sprintf(msg, fields...))
}

// Printf interface for kafka's implementation.
func (l *Zap) Printf(msg string, fields ...interface{}) {
	l.Infof(msg, fields...)
}

// initConfig method that initializes logger.
func initConfig(level string) (zap.Config, error) {
	cfg := zap.NewProductionConfig()

	lvl, err := zap.ParseAtomicLevel(level)
	if err != nil {
		emptyConfig := zap.Config{}
		return emptyConfig, fmt.Errorf("init Zap config: %w", err)
	}
	cfg.Level = lvl
	cfg.DisableCaller = true

	return cfg, nil
}

// LogHandler handler for requests logging.
func (l *Zap) LogHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		loggingWriter := createResponseWriter(w)
		h.ServeHTTP(loggingWriter, r)
		duration := time.Since(start)

		l.zap.Info("new incoming HTTP request",
			zap.String("uri", r.RequestURI),
			zap.String("method", r.Method),
			zap.Int("status", loggingWriter.getResponseData().status),
			zap.String("content-type", loggingWriter.Header().Get("Content-Type")),
			zap.String("content-encoding", loggingWriter.Header().Get("Content-Encoding")),
			zap.String("HashSHA256", loggingWriter.Header().Get("HashSHA256")),
			zap.Duration("duration", duration),
			zap.Int("size", loggingWriter.getResponseData().size),
		)
	})
}
