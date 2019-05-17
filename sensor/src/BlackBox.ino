#include "Adafruit_BME680.h"

#define SEALEVELPRESSURE_HPA (1013.25)


Adafruit_BME680 bme; // I2C

const int ledPin = D6;
const uint32_t CATASTROPHE_LENGTH_SECONDS = 5 * 60;
const uint32_t CATASTROPHE_HOLD_LENGTH_SECONDS = 5 * 60;
uint32_t catastropheTimer = 0;
uint32_t catastropheHoldTimer = 0;
const uint32_t SENSOR_DELAY = 10;
uint32_t sensorTimer = 0;
char deviceName[32] = "unknown";

const double TEMPERATURE_MAX = 600;
const double HUMIDITY_MIN = 20;
const double PRESSURE_MIN = 949;
const double AIR_QUALITY_MAX = 500;

// Define methods
int introduceCatastrophe(String cmd);
double getCatastrophePercent();
double calculateValue(double initial, double target, double percent);
void deviceNameHandler(const char *topic, const char *data);
void readSensor();
void blink(unsigned long delayMs);

/*
 * Setup the device at startup.
 */
void setup() {
  // Register the catastrophe function
  Particle.function("catastrophe", introduceCatastrophe);

  // Get the device name to include in the event data
  Particle.subscribe("particle/device/name", deviceNameHandler);
  Particle.publish("particle/device/name");

  // Setup pin mode for the LED
  pinMode(ledPin, OUTPUT);

  // Setup BME Sensor
  if (!bme.begin()) {
    Particle.publish("Log", "Could not find a valid BME680 sensor, check wiring!", PRIVATE);
  } else {
    Particle.publish("Log", "bme.begin() success =)", PRIVATE);
    // Set up oversampling and filter initialization
    bme.setTemperatureOversampling(BME680_OS_8X);
    bme.setHumidityOversampling(BME680_OS_2X);
    bme.setPressureOversampling(BME680_OS_4X);
    bme.setIIRFilterSize(BME680_FILTER_SIZE_3);
    bme.setGasHeater(320, 150); // 320*C for 150 ms
    sensorTimer = Time.now();
  }
}

/*
 * Main loop
 */
void loop() {
  if (Time.now() > sensorTimer) {
    readSensor();
    sensorTimer = Time.now() + SENSOR_DELAY;
  } else {
    blink(250);
  }

  // Delay 1 second
  delay(1000);
}

/*
 * Read the sensor data and publish to the Particle Cloud.
 */
void readSensor() {
  if (! bme.performReading()) {
    Particle.publish("Log", "Failed to perform reading :(", PRIVATE);
    return;
  }

  // Wait till we have a name
  if (deviceName == "unknown") {
    return;
  }

  // Turn on LED
  digitalWrite(ledPin, HIGH);

  // Read the sensor
  double temperatureInC = bme.temperature;
  double relativeHumidity = bme.humidity;
  double pressureHpa = bme.pressure / 100.0;
  double gasResistanceKOhms = bme.gas_resistance / 1000.0;

  // Modify data based on catastrophe
  double percent = getCatastrophePercent();
  if (percent > 0.01) {
    temperatureInC = calculateValue(temperatureInC, TEMPERATURE_MAX, percent);
    relativeHumidity = calculateValue(relativeHumidity, HUMIDITY_MIN, percent);
    pressureHpa = calculateValue(pressureHpa, PRESSURE_MIN, percent);
    gasResistanceKOhms = calculateValue(gasResistanceKOhms, AIR_QUALITY_MAX, percent);
  }

  // Publish data
  String data = String::format(
    "{\"temperature\":%.2f, \"humidity\":%.2f, \"pressure\":%.2f, \"airQuality\":%.2f, \"device\":\"%s\"}",
    temperatureInC, relativeHumidity, pressureHpa, gasResistanceKOhms, deviceName
  );
  Particle.publish("sensor-data", data, PRIVATE);
  Particle.publish("catastrophe-timer", String(catastropheTimer), PRIVATE);

  // Blink rapidly
  blink(125);
  blink(125);
  blink(125);
}

/*
 * Callback for grabbing the device name.
 */
void deviceNameHandler(const char *topic, const char *data) {
  strncpy(deviceName, data, sizeof(deviceName) - 1);
}

/*
 * Returns the current multiplier for a catastrophe.
 *
 * Returns 0 if there is no catastrophe.
 */
double getCatastrophePercent() {
  // Always 100 percent while holding;
  if (catastropheHoldTimer > 0) {
    uint32_t currentTime = Time.now();
    if (currentTime > catastropheHoldTimer) {
      catastropheHoldTimer = 0;
    }
    return 1.0;
  }

  // Now check if ramping up to a catastrophe
  if (catastropheTimer == 0) {
    return 0;
  }

  double percent = 0.0;
  uint32_t currentTime = Time.now();
  if (currentTime < catastropheTimer) {
      uint32_t delta = catastropheTimer - currentTime;
      percent = 1 - ((double) delta / (double) CATASTROPHE_LENGTH_SECONDS);
  } else {
      // Catastrophe over, reset the timer
      catastropheTimer = 0;
      catastropheHoldTimer = currentTime + CATASTROPHE_HOLD_LENGTH_SECONDS;
  }
  return percent;
}

/*
 * Calculates the value based on the catastrophe multipler.
 */
double calculateValue(double initial, double target, double percent) {
  if (initial > target) {
    return initial - (percent * (initial - target));
  } else {
    return initial + (percent * (target - initial));
  }
}

/*
 * Callback for a Particle Function to initate a catastrophe.
 */
int introduceCatastrophe(String cmd) {
  // Don't reset the timer if already active
  if (catastropheTimer > 0) {
    return 0;
  }

  catastropheTimer = Time.now() + CATASTROPHE_LENGTH_SECONDS;
  return 1;
}

/*
 * Blinks the LED once with the given delay.
 */
void blink(unsigned long delayMs) {
  digitalWrite(ledPin, HIGH);
  delay(delayMs);
  digitalWrite(ledPin, LOW);
  delay(delayMs);
}
