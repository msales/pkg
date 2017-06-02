package log

var (
	// Null is the null Logger instance.
	Null = &nullLogger{}
)

// Logger represents an abstract logging object.
type Logger interface {
	// Debug logs a debug message.
	Debug(msg string, ctx ...interface{})
	// Info logs an informational message.
	Info(msg string, ctx ...interface{})
	// Error logs an error message.
	Error(msg string, ctx ...interface{})
}

type nullLogger struct{}

// Debug logs a debug message.
func (l nullLogger) Debug(msg string, ctx ...interface{}) {}

// Info logs an informational message.
func (l nullLogger) Info(msg string, ctx ...interface{}) {}

// Error logs an error message.
func (l nullLogger) Error(msg string, ctx ...interface{}) {}
