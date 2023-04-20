// Author: Chuck Sailer
// Date: 2023-04-20
// Description: This file contains the code for the CloudWatchLogger struct and its methods.
package logging

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
)

// CloudWatchLogger is a logger that writes log messages to CloudWatch.
type CloudWatchLogger struct {
	svc           *cloudwatchlogs.CloudWatchLogs
	logGroupName  string
	logStreamName string
	logChan       chan LogMessage
	doneChan      chan bool
}

// LogMessage is a struct that contains the information needed to write a log message to CloudWatch.
type LogMessage struct {
	LogGroupName  string // Name of the log group to write to
	LogStreamName string // Name of the log stream to write to
	Message       string // Log message to write
}

// NewCloudWatchLogger creates a new CloudWatchLogger.
func NewCloudWatchLogger(logGroupName, logStreamName string) (*CloudWatchLogger, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
	})
	if err != nil {
		return nil, err
	}

	cwl := &CloudWatchLogger{
		svc:           cloudwatchlogs.New(sess),
		logGroupName:  logGroupName,
		logStreamName: logStreamName,
		logChan:       make(chan LogMessage, 100),
		doneChan:      make(chan bool),
	}

	go cwl.writeToCloudWatch(cwl.logChan, cwl.doneChan)

	return cwl, nil
}

// NewLogMessage creates a new LogMessage struct.
func NewLogMessage(logGroupName, logStreamName, message string) LogMessage {
	return LogMessage{
		LogGroupName:  logGroupName,
		LogStreamName: logStreamName,
		Message:       message,
	}
}

// WriteToCloudWatchLog uses Go's channels to handle incoming log messages asynchronously.
// The function takes two channels as parameters:
// - logChan: a channel that contains log messages to be written to CloudWatch
// - doneChan: a channel that is used to signal the Write function to stop
//
// Example:
//
//	logChan := make(chan LogMessage, 100)  //-- create a buffered channel to hold up to 100 log messages
//	doneChan := make(chan bool)
//	go cwl.WriteToCloudWatchLog(logChan, doneChan)
//
//	Once you've set up the channels and started the logger (go cwl.WriteToCloudWatchLog()), you can write log messages
//	to the logChan channel using a LogMessage struct:
//
//	logMsg := LogMessage{
//		LogGroupName:  "my-log-group",
//		LogStreamName: "my-log-stream",
//		Message:       "this is a test log message",
//	}
//
//	logChan <- logMsg
func (cwl *CloudWatchLogger) writeToCloudWatch(logChan chan LogMessage, doneChan chan bool) error {
	// Loop over incoming log messages and write them to CloudWatch
	for logMsg := range logChan {
		_, err := cwl.svc.PutLogEvents(&cloudwatchlogs.PutLogEventsInput{
			LogGroupName:  aws.String(logMsg.LogGroupName),
			LogStreamName: aws.String(logMsg.LogStreamName),
			LogEvents: []*cloudwatchlogs.InputLogEvent{
				{
					Timestamp: aws.Int64(time.Now().UnixNano() / int64(time.Millisecond)),
					Message:   aws.String(logMsg.Message),
				},
			},
		})
		if err != nil {
			fmt.Printf("failed to put CloudWatch log event for log group %s and log stream %s: %v\n", logMsg.LogGroupName, logMsg.LogStreamName, err)
		}
	}

	// Signal that the logger has finished processing messages
	doneChan <- true
	return nil
}

// Debug writes a debug log message to CloudWatch.
func (cwl *CloudWatchLogger) Debug(msg string, args ...interface{}) {
	logMsg := NewLogMessage(cwl.logGroupName, cwl.logStreamName, fmt.Sprintf(msg, args...))
	cwl.logChan <- logMsg
}

// Info writes an informational log message to CloudWatch.
func (cwl *CloudWatchLogger) Info(msg string, args ...interface{}) {
	logMsg := NewLogMessage(cwl.logGroupName, cwl.logStreamName, fmt.Sprintf(msg, args...))
	cwl.logChan <- logMsg
}

// Warn writes a warning log message to CloudWatch.
func (cwl *CloudWatchLogger) Warn(msg string, args ...interface{}) {
	logMsg := NewLogMessage(cwl.logGroupName, cwl.logStreamName, fmt.Sprintf(msg, args...))
	cwl.logChan <- logMsg
}

// Error writes an error log message to CloudWatch.
func (cwl *CloudWatchLogger) Error(msg string, args ...interface{}) {
	logMsg := NewLogMessage(cwl.logGroupName, cwl.logStreamName, fmt.Sprintf(msg, args...))
	cwl.logChan <- logMsg
}

// Close stops the CloudWatch logger.
func (cwl *CloudWatchLogger) Close() error {
	close(cwl.logChan)
	<-cwl.doneChan // Wait for the logger to finish processing messages
	return nil
}
