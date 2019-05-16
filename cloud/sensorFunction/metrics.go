package blackbox

import (
	"context"
	"fmt"
	"time"

	"contrib.go.opencensus.io/exporter/stackdriver"
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"
)

var (
	DeviceKey, _ = tag.NewKey("device")
	TemperatureM    = stats.Float64("blackbox/m/temperature/last", "The temperature in celsius", "c")
	HumidityM      = stats.Float64("blackbox/m/humidity/last", "The humidity in relative humidity", "rh")
	PressureM      = stats.Float64("blackbox/m/pressure/last", "The pressure in hPa", "hPa")
	AirQualityM    = stats.Float64("blackbox/m/air_quality/last", "The indoor air quality using IAQ index", "iaq")
	TemperatorView = &view.View{
		Name:        "blackbox/temperature/last",
		Measure:     TemperatureM,
		Description: "The distribution of the temperatures",
		Aggregation: view.LastValue(),
		TagKeys:     []tag.Key{DeviceKey},
	}
	HumidityView   =&view.View{
		Name:        "blackbox/humidity/last",
		Measure:     HumidityM,
		Description: "The distribution of the humidity",
		Aggregation: view.LastValue(),
		TagKeys:     []tag.Key{DeviceKey},
	}
	PressuareView  = &view.View{
		Name:        "blackbox/pressure/last",
		Measure:     PressureM,
		Description: "The distribution of the pressure",
		Aggregation: view.LastValue(),
		TagKeys:     []tag.Key{DeviceKey},
	}
	AirQualityView = &view.View{
		Name:        "blackbox/air_quality/last",
		Measure:     AirQualityM,
		Description: "The distribution of the air quality",
		Aggregation: view.LastValue(),
		TagKeys:     []tag.Key{DeviceKey},
	}
	Stackdriver *stackdriver.Exporter
	metricsRegistered = false
)

// NewMetrics returns a newly created Metrics struct.
func RegisterMetrics() (err error) {
	// Register the Stackdriver exporter
	if metricsRegistered {
		return nil
	}
	err = view.Register(TemperatorView, HumidityView, PressuareView, AirQualityView)
	if err != nil {
		return fmt.Errorf("Failed to register the views: %v", err)
	}
	Stackdriver, err = stackdriver.NewExporter(stackdriver.Options{})
	if err != nil {
		return fmt.Errorf("Failed to create Stackdriver exporter: %v", err)
	}
	if err := Stackdriver.StartMetricsExporter(); err != nil {
		return fmt.Errorf("Error starting metric exporter: %v", err)
	}

	// Register it as a metrics exporter
	view.RegisterExporter(Stackdriver)
	view.SetReportingPeriod(60 * time.Second)
	metricsRegistered = true

	return nil
}

// RecordEvent records the event to StackDriver
// Only expecting to call this once per Cloud Function invocation, so flushing
// at the end of the function.
func RecordEvent(ctx context.Context, event *SensorEvent) (err error) {
	ctx, err = tag.New(ctx, tag.Insert(DeviceKey, event.Device))
	if err != nil {
		return err
	}

	stats.Record(ctx,
		TemperatureM.M(event.Temperature),
		HumidityM.M(event.Humidity),
		PressureM.M(event.Pressure),
		AirQualityM.M(event.AirQuality),
	)
	return nil
}
