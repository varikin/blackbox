// Package blackbox provides a Cloud function to send
// sensor data to a Googl Sheet
package blackbox

import (
	"context"
	"encoding/json"
	"fmt"
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

// Run consumes a Pub/Sub message.
func Run(ctx context.Context, m PubSubMessage) error {
	logger := NewCloudFunctionLogger()
	event, err := DecodeSensorEvent(m.Data)
	if err != nil {
		logger.error("Failed to get sensor data: %s", err.Error())
		return err
	}
	logger.log("Got sensor data %+v", event)

	return nil
}
