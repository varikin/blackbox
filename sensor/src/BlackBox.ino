/***************************************************************************
  This is a library for the BME680 gas, humidity, temperature & pressure sensor

  Designed specifically to work with the Adafruit BME680 Breakout
  ----> http://www.adafruit.com/products/3660

  These sensors use I2C or SPI to communicate, 2 or 4 pins are required
  to interface.

  Adafruit invests time and resources providing this open source code,
  please support Adafruit and open-source hardware by purchasing products
  from Adafruit!

  Written by Limor Fried & Kevin Townsend for Adafruit Industries.
  BSD license, all text above must be included in any redistribution
 ***************************************************************************/

#include "Adafruit_BME680.h"

#define SEALEVELPRESSURE_HPA (1013.25)

Adafruit_BME680 bme; // I2C

double temperatureInC = 0;
double relativeHumidity = 0;
double pressureHpa = 0;
double gasResistanceKOhms = 0;

const uint32_t CATASTROPHE_LENGTH = 120;
uint32_t catastropheTimer = 0;
int introduceCatastrophe(String cmd);
double getMultiplier();
double calculateCatastrophe(double multiplier, double value);

void setup() {
  Particle.function("catastrophe", introduceCatastrophe);
  if (!bme.begin()) {
    Particle.publish("Log", "Could not find a valid BME680 sensor, check wiring!");
  } else {
    Particle.publish("Log", "bme.begin() success =)");
    // Set up oversampling and filter initialization
    bme.setTemperatureOversampling(BME680_OS_8X);
    bme.setHumidityOversampling(BME680_OS_2X);
    bme.setPressureOversampling(BME680_OS_4X);
    bme.setIIRFilterSize(BME680_FILTER_SIZE_3);
    bme.setGasHeater(320, 150); // 320*C for 150 ms

    Particle.variable("temperature", &temperatureInC, DOUBLE);
    Particle.variable("humidity", &relativeHumidity, DOUBLE);
    Particle.variable("pressure", &pressureHpa, DOUBLE);
    Particle.variable("gas", &gasResistanceKOhms, DOUBLE);
  }
}

void loop() {
  if (! bme.performReading()) {
    Particle.publish("Log", "Failed to perform reading :(");
  } else {
    temperatureInC = bme.temperature;
    relativeHumidity = bme.humidity;
    pressureHpa = bme.pressure / 100.0;
    gasResistanceKOhms = bme.gas_resistance / 1000.0;

    double multiplier = getMultiplier();
    if (multiplier > 0.0) {
      temperatureInC = calculateCatastrophe(multiplier, temperatureInC);
      relativeHumidity = calculateCatastrophe(multiplier, relativeHumidity);
      pressureHpa = calculateCatastrophe(multiplier, pressureHpa);
      gasResistanceKOhms = calculateCatastrophe(multiplier, gasResistanceKOhms);
    }

    String data = String::format(
      "{"
        "\"temperature\":%.2f,"
        "\"humidity\":%.2f,"
        "\"pressure\":%.2f,"
        "\"airQuality\":%.2f"
      "}",
      temperatureInC,
      relativeHumidity,
      pressureHpa,
      gasResistanceKOhms
    );

    Particle.publish("sensor-data", data, 60, PRIVATE);
    Particle.publish("catastrophe-timer", String(catastropheTimer), 60, PRIVATE);
  }
  delay(10 * 1000);
}

double getMultiplier() {
  if (catastropheTimer == 0) {
    return 0;
  }

  double multiplier = 0.0;
  uint32_t currentTime = Time.now();
  if (currentTime < catastropheTimer) {
      uint32_t delta = catastropheTimer - currentTime;
      multiplier = (double) delta / (double) CATASTROPHE_LENGTH;
  } else {
      // Catastrophe over, reset the timer
      catastropheTimer = 0;
  }
  return multiplier;
}

double calculateCastrophe(double multipler, double value) {
  return value + (multipler * value);
}

int introduceCatastrophe(String cmd) {
  // Don't reset the timer if already active
  if (catastropheTimer > 0) {
    return 0;
  }

  catastropheTimer = Time.now() + CATASTROPHE_LENGTH;
  return 1;
}
