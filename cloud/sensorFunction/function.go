// Package blackbox provides a Cloud function to send
// sensor data to a Googl Sheet
package blackbox

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"
)

// PubSubMessage is the payload of a Pub/Sub event. Please refer to the docs for
// additional information regarding Pub/Sub events.
type PubSubMessage struct {
	Data        []byte            `json:"data"`
	Attributes  map[string]string `json:"attributes"`
	PublishTime string            `json:"publishTime"`
}

// Run consumes a Pub/Sub message.
func Run(ctx context.Context, m PubSubMessage) error {
	logger := NewCloudFunctionLogger()
	err := RegisterMetrics()
	if err != nil {
		logger.error(err.Error())
		return err
	}

	timestamp, err := getTimestamp(m.Attributes)
	if err != nil {
		logger.error(err.Error())
		return err
	}

	// Get the sheet id from the environment
	sheetID := os.Getenv("SHEET_ID")
	if len(sheetID) == 0 {
		err := errors.New("SHEET_ID environment variable not set")
		logger.error(err.Error())
		return err
	}

	// Read the event off the data
	event, err := DecodeSensorEvent(m.Data, timestamp)
	if err != nil {
		logger.error("Failed to get sensor data: %s", err.Error())
		return err
	}

	// Record the even to StackDriver
	if err = RecordEvent(ctx, event); err != nil {
		logger.error("Failed to append data to the sheet: error=%s", err.Error())
		return err
	}

	// Append the event to sheet
	if err = AppendEventToSpreadsheet(ctx, sheetID, event); err != nil {
		logger.error("Failed to append data to the sheet: error=%s", err.Error())
		return err
	}
	logger.log("Reported event: %v+", event)

	return nil
}

func getTimestamp(attributes map[string]string) (string, error) {
	timestamp, err := time.Parse("2006-01-02T15:04:05.999Z", attributes["published_at"])
	if err != nil {
		return "", fmt.Errorf("Failed to parse the timestamp: published_at=%s, error=%s", attributes["published_at"], err.Error())
	}
	return timestamp.Format("01/02/2006 3:03 PM"), nil
}
