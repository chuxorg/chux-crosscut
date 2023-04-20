package logging

import "testing"

func TestLoggerFacade(t *testing.T) {
	// Set up the CloudWatchLogger and FileLogger
	cwl, err := NewCloudWatchLogger("my-log-group", "my-log-stream")
	if err != nil {
		t.Fatalf("Error creating CloudWatch logger: %v", err)
	}
	fl, err := NewFileLogger("./logs", "test", 1024)
	if err != nil {
		t.Fatalf("Error creating File logger: %v", err)
	}

	// Set up the logger facade with the CloudWatchLogger and FileLogger
	loggerFacade := LoggerFacade{
		logger: NewLoggerComposition(cwl, fl),
	}

	// Write some log messages using the logger facade
	loggerFacade.Debug("This is a debug message")
	loggerFacade.Info("This is an informational message")
	loggerFacade.Warn("This is a warning message")
	loggerFacade.Error("This is an error message")

	// Close the logger facade
	err = loggerFacade.logger.Close()
	if err != nil {
		t.Errorf("Error closing logger facade: %v", err)
	}
}
