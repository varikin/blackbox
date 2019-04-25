/*
 * Checks the state of a momentary switch and publishes to Particle when pressed.
 */

int momPin = D5;
int previousMomState = 0;

// Setup the device with the momentary switch pin in the right mode.
void setup() {
  pinMode(momPin, INPUT_PULLDOWN);
  Serial.begin(115200);
}

// Main loop, check the switch state
void loop() {

  int momState = digitalRead(momPin);

  // Only act if the state changed from the last run.
  if (previousMomState != momState) {
    Serial.printf("Arcade: %d\n", momState);
    previousMomState = momState;
    if (momState == 1) {
      Particle.publish("simulate-catastrophe", String(momState), 60, PRIVATE);
    }
  }

  // Might miss a state change if someone presses the switch and releases within 100ms
  // But also, seems excessive to not sleep between runs.
  // Realistically, who will press the button that fast?
  delay(100);
}

