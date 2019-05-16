package blackbox

import (
	"context"
	"fmt"

	"google.golang.org/api/sheets/v4"
)

// AppendEventToSpreadsheet appends a SensorEvent to a Google Sheet.
func AppendEventToSpreadsheet(ctx context.Context, sheetID string, event *SensorEvent) error {
	sheetService, err := sheets.NewService(ctx)
	if err != nil {
		return fmt.Errorf("Failed to get Sheets client: error=%s", err.Error())
	}

	values := []interface{}{
		event.Timestamp,
		event.Device,
		event.Temperature,
		event.Humidity,
		event.Pressure,
		event.AirQuality,
	}
	var valueRange sheets.ValueRange
	valueRange.Values = append(valueRange.Values, values)
	_, err = sheetService.Spreadsheets.Values.Append(sheetID, "A1", &valueRange).ValueInputOption("USER_ENTERED").Do()
	if err != nil {
		return fmt.Errorf("Failed to append data to sheet: data=%+v, error=%s", event, err.Error())
	}
	return nil
}
