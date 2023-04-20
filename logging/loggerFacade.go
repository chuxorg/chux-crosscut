package logging

import "errors"

// Configuration for determining which logger to use
type LoggerConfig struct {
	UseCloudWatch bool   // Whether to use CloudWatchLogger or not
	UseFileLogger bool   // Whether to use FileLogger or not
	LogDirectory  string // Directory to store log files in for FileLogger
	LogFilePrefix string // Prefix to use for log files for FileLogger
	MaxFileSize   int64  // Maximum size for each log file for FileLogger
	LogGroupName  string // Log group name for CloudWatchLogger
	LogStreamName string // Log stream name for CloudWatchLogger
}

type LoggerFacade struct {
	logger ILogger // Underlying logger
}

// Creates a new LoggerFacade based on the specified configuration
func NewLoggerFacade(config LoggerConfig) (*LoggerFacade, error) {
	// Check that either UseCloudWatch or UseFileLogger is set to true
	if !config.UseCloudWatch && !config.UseFileLogger {
		return nil, errors.New("either UseCloudWatch or UseFileLogger must be true")
	}

	// Check that only one of UseCloudWatch or UseFileLogger is set to true
	if config.UseCloudWatch && config.UseFileLogger {
		return nil, errors.New("only one of UseCloudWatch or UseFileLogger can be true")
	}

	// Create and return a new LoggerFacade based on the specified configuration
	if config.UseCloudWatch {
		cwLogger, err := NewCloudWatchLogger(config.LogGroupName, config.LogStreamName)
		if err != nil {
			return nil, err
		}
		return &LoggerFacade{logger: cwLogger}, nil
	} else {
		fileLogger, err := NewFileLogger(config.LogDirectory, config.LogFilePrefix, config.MaxFileSize)
		if err != nil {
			return nil, err
		}
		return &LoggerFacade{logger: fileLogger}, nil
	}
}

// Implements the ILogger interface by forwarding log messages to the underlying logger
func (lf *LoggerFacade) Debug(msg string, args ...interface{}) {
	lf.logger.Debug(msg, args...)
}

func (lf *LoggerFacade) Info(msg string, args ...interface{}) {
	lf.logger.Info(msg, args...)
}

func (lf *LoggerFacade) Warn(msg string, args ...interface{}) {
	lf.logger.Warn(msg, args...)
}

func (lf *LoggerFacade) Error(msg string, args ...interface{}) {
	lf.logger.Error(msg, args...)
}

func (lf *LoggerFacade) Close() {
	lf.logger.Close()
}
