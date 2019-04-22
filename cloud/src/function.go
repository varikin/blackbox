// Package blackbox provides a Cloud function to send
// sensor data to a Googl Sheet
package blackbox

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"google.golang.org/api/sheets/v4"
)

// PubSubMessage is the payload of a Pub/Sub event. Please refer to the docs for
// additional information regarding Pub/Sub events.
type PubSubMessage struct {
	Data []byte `json:"data"`
}

// SensorEvent is the payload within the Pub/Sub event data.
type SensorEvent struct {
	Temperature float64 `json:"temperature"`
	Humidity    float64 `json:"humidity"`
	Pressure    float64 `json:"pressure"`
	AirQuality  float64 `json:"airQuality"`
}

// DecodeSensorEvent decodes the raw sensor event and converts to a SensorEvent struct
func DecodeSensorEvent(data []byte) (*SensorEvent, error) {
	var event SensorEvent
	if err := json.Unmarshal(data, &event); err != nil {
		return nil, fmt.Errorf("Failed to unmarshal sensor data: data=%s, error=%s", string(data), err.Error())
	}
	return &event, nil
}

// AppendEvent appends a SensorEvent to a Google Sheet.
func AppendEvent(ctx context.Context, sheetID string, event *SensorEvent) error {
	sheetService, err := sheets.NewService(ctx)
	if err != nil {
		return fmt.Errorf("Failed to get Sheets client: error=%s", err.Error())
	}

	values := []interface{}{event.Temperature, event.Humidity, event.Pressure, event.AirQuality}
	var valueRange sheets.ValueRange
	valueRange.Values = append(valueRange.Values, values)
	_, err = sheetService.Spreadsheets.Values.Append(sheetID, "A1", &valueRange).ValueInputOption("USER_ENTERED").Do()
	if err != nil {
		return fmt.Errorf("Failed to append data to sheet: data=%+v, error=%s", event, err.Error())
	}
	return nil

}

// Run consumes a Pub/Sub message.
func Run(ctx context.Context, m PubSubMessage) error {
	logger := NewCloudFunctionLogger()

	// Get the sheet id from the environment
	sheetID := os.Getenv("SHEET_ID")
	if len(sheetID) == 0 {
		err := errors.New("SHEET_ID environment variable not set")
		logger.error(err.Error())
		return err
	}

	// Read the event off the data
	event, err := DecodeSensorEvent(m.Data)
	if err != nil {
		logger.error("Failed to get sensor data: %s", err.Error())
		return err
	}

	// Append the event to sheet
	if err = AppendEvent(ctx, sheetID, event); err != nil {
		logger.error("Failed to append data to the sheet: error=%s", err.Error())
		return err
	}
	logger.log("Appended data to the sheet")

	return nil
}
