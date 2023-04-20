package logging

import "fmt"

// LoggerComposition is a composition of multiple loggers.
type LoggerComposition struct {
	loggers []ILogger
}

// NewLoggerComposition creates a new logger composition with the given loggers.
func NewLoggerComposition(loggers ...ILogger) ILogger {
	return &LoggerComposition{loggers: loggers}
}

// Debug writes a debug log message to all loggers in the composition.
func (lc *LoggerComposition) Debug(msg string, args ...interface{}) {
	for _, logger := range lc.loggers {
		logger.Debug(msg, args...)
	}
}

// Info writes an informational log message to all loggers in the composition.
func (lc *LoggerComposition) Info(msg string, args ...interface{}) {
	for _, logger := range lc.loggers {
		logger.Info(msg, args...)
	}
}

// Warn writes a warning log message to all loggers in the composition.
func (lc *LoggerComposition) Warn(msg string, args ...interface{}) {
	for _, logger := range lc.loggers {
		logger.Warn(msg, args...)
	}
}

// Error writes an error log message to all loggers in the composition.
func (lc *LoggerComposition) Error(msg string, args ...interface{}) {
	for _, logger := range lc.loggers {
		logger.Error(msg, args...)
	}
}

// Close closes all loggers in the composition.
func (lc *LoggerComposition) Close() error {
	var errs []error
	for _, logger := range lc.loggers {
		err := logger.Close()
		if err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf("error closing loggers: %v", errs)
	}
	return nil
}
