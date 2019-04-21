package blackbox

import (
	"log"
	"os"
)

// SimpleLogger is a simple logger
type SimpleLogger interface {
	log(v ...interface{})
	error(v ...interface{})
}

// CloudFunctionLogger manages a stdout and stderr logger for Google Cloud Functions
type CloudFunctionLogger struct {
	logger      *log.Logger
	errorLogger *log.Logger
}

// NewCloudFunctionLogger returns a new CloudFunctionLogger
func NewCloudFunctionLogger() *CloudFunctionLogger {
	return &CloudFunctionLogger{
		logger:      log.New(os.Stdout, "", 0),
		errorLogger: log.New(os.Stderr, "", 0),
	}
}

func (cfl *CloudFunctionLogger) log(v ...interface{}) {
	cfl.logger.Println(v)
}

func (cfl *CloudFunctionLogger) error(v ...interface{}) {
	cfl.errorLogger.Println(v)
}
