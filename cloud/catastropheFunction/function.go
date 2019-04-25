// Package blackbox provides a Cloud function to send
// sensor data to a Googl Sheet
package blackbox

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

// Run is the entry point into the Function.
// The pub/sub message is not needed, it is just a trigger.
func Run(ctx context.Context, _ interface{}) error {
	logger := NewCloudFunctionLogger()

	// Get the access token for Particle
	accessToken := os.Getenv("PARTICLE_ACCESS_TOKEN")
	if len(accessToken) == 0 {
		err := errors.New("PARTICLE_ACCESS_TOKEN environment variable not set")
		logger.error(err.Error())
		return err
	}

	// Get the list of devices that could be acted upon
	rawDeviceIDs := os.Getenv("DEVICE_IDS")
	if len(rawDeviceIDs) == 0 {
		err := errors.New("DEVICE_IDS environment variable not set")
		logger.error(err.Error())
		return err
	}
	deviceIDs := strings.Split(rawDeviceIDs, ",")

	// Select one device randomly
	seed := rand.NewSource(time.Now().UnixNano())
	random := rand.New(seed)
	index := random.Intn(len(deviceIDs))
	deviceID := deviceIDs[index]

	// Post to the Particle Cloud for the device
	api := fmt.Sprintf("https://api.particle.io/v1/devices/%s/catastrophe", deviceID)
	_, err := http.PostForm(api, url.Values{
		"access_token": {accessToken},
	})

	if err != nil {
		err := fmt.Errorf("Error simulating catastrophe: deviceId=%s, error=%s", deviceID, err.Error())
		logger.error(err.Error())
		return err
	}
	logger.log("Simulated catastrophe on device %s", deviceID)
	return nil
}
