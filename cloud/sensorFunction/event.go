package blackbox

import (
	"encoding/json"
	"fmt"
)

// SensorEvent is the payload within the Pub/Sub event data.
type SensorEvent struct {
	Temperature float64 `json:"temperature"`
	Humidity    float64 `json:"humidity"`
	Pressure    float64 `json:"pressure"`
	AirQuality  float64 `json:"airQuality"`
	Device      string  `json:"device"`
	Timestamp   string  `json:"timestamp"`
}

// DecodeSensorEvent decodes the raw sensor event and converts to a SensorEvent struct
func DecodeSensorEvent(data []byte, timestamp string) (*SensorEvent, error) {
	var event SensorEvent
	if err := json.Unmarshal(data, &event); err != nil {
		return nil, fmt.Errorf("Failed to unmarshal sensor data: data=%s, error=%s", string(data), err.Error())
	}
	event.Timestamp = timestamp
	return &event, nil
}
