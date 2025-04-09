package logger

import (
	"os"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	// Global logger instance
	globalLogger *zap.SugaredLogger
	once         sync.Once
)

// Config holds configuration for the logger
type Config struct {
	// Development puts the logger in development mode, which changes the
	// behavior of DPanicLevel and makes it log at the stack trace of all error logs.
	Development bool
	// Level is the minimum enabled logging level.
	Level string
	// OutputPaths is a list of URLs or file paths to write logging output to.
	// Default is "stdout".
	OutputPaths []string
	// ErrorOutputPaths is a list of URLs or file paths to write internal logger errors to.
	// Default is "stderr".
	ErrorOutputPaths []string
	// Format specifies the output format: "json" or "console"
	Format string
}

// DefaultConfig returns the default logger configuration
func DefaultConfig() Config {
	return Config{
		Development:      false,
		Level:            "info",
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
		Format:           "console", // Default to console format
	}
}

// Init initializes the global logger with the given configuration
func Init(cfg Config) {
	once.Do(func() {
		// Parse log level
		var level zapcore.Level
		if err := level.UnmarshalText([]byte(cfg.Level)); err != nil {
			level = zapcore.InfoLevel
		}

		// Create encoder config
		encoderConfig := zapcore.EncoderConfig{
			TimeKey:        "ts",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			FunctionKey:    zapcore.OmitKey,
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.CapitalLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		}

		// Determine encoding format
		encoding := "console"
		if cfg.Format == "json" {
			encoding = "json"
		}

		// For console format, customize the encoder config for better readability
		if encoding == "console" {
			encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
			encoderConfig.ConsoleSeparator = " "
		}

		// Create logger config
		zapConfig := zap.Config{
			Level:             zap.NewAtomicLevelAt(level),
			Development:       cfg.Development,
			DisableCaller:     false,
			DisableStacktrace: !cfg.Development,
			Sampling:          nil,
			Encoding:          encoding,
			EncoderConfig:     encoderConfig,
			OutputPaths:       cfg.OutputPaths,
			ErrorOutputPaths:  cfg.ErrorOutputPaths,
		}

		// Build logger
		logger, err := zapConfig.Build(zap.AddCallerSkip(0))
		if err != nil {
			// If we can't build the logger, fall back to a basic logger
			logger = zap.New(zapcore.NewCore(
				zapcore.NewConsoleEncoder(encoderConfig),
				zapcore.AddSync(os.Stdout),
				zapcore.InfoLevel,
			))
		}

		// Create sugared logger
		globalLogger = logger.Sugar()
	})
}

// Get returns the global logger instance
// If the logger hasn't been initialized, it initializes it with default config
func Get() *zap.SugaredLogger {
	if globalLogger == nil {
		Init(DefaultConfig())
	}
	return globalLogger
}

// Debug logs a message at debug level
func Debug(args ...interface{}) {
	withCaller().Debug(args...)
}

// Debugf logs a formatted message at debug level
func Debugf(format string, args ...interface{}) {
	withCaller().Debugf(format, args...)
}

// Info logs a message at info level
func Info(args ...interface{}) {
	withCaller().Info(args...)
}

// Infof logs a formatted message at info level
func Infof(format string, args ...interface{}) {
	withCaller().Infof(format, args...)
}

// Warn logs a message at warn level
func Warn(args ...interface{}) {
	withCaller().Warn(args...)
}

// Warnf logs a formatted message at warn level
func Warnf(format string, args ...interface{}) {
	withCaller().Warnf(format, args...)
}

// Error logs a message at error level
func Error(args ...interface{}) {
	withCaller().Error(args...)
}

// Errorf logs a formatted message at error level
func Errorf(format string, args ...interface{}) {
	withCaller().Errorf(format, args...)
}

// Fatal logs a message at fatal level and then calls os.Exit(1)
func Fatal(args ...interface{}) {
	withCaller().Fatal(args...)
}

// Fatalf logs a formatted message at fatal level and then calls os.Exit(1)
func Fatalf(format string, args ...interface{}) {
	withCaller().Fatalf(format, args...)
}

// With returns a logger with the specified key-value pairs
func With(args ...interface{}) *zap.SugaredLogger {
	return withCaller().With(args...)
}

// Sync flushes any buffered log entries
func Sync() error {
	return Get().Sync()
}

// withCaller returns a logger with the caller skip level increased
// to skip the wrapper functions in this package
func withCaller() *zap.SugaredLogger {
	// Get the underlying zap.Logger
	zapLogger := Get().Desugar()

	// Add a caller skip to bypass the wrapper function
	zapLogger = zapLogger.WithOptions(zap.AddCallerSkip(1))

	// Return a sugared logger
	return zapLogger.Sugar()
}
