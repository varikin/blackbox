/*
 * Checks the state of a momentary switch and publishes to Particle when pressed.
 */
const int LED_PIN = D3;
const int MOM_PIN = D5;
int previousMomState = 0;

void blink(unsigned long delayMs);

// Setup the device with the momentary switch pin in the right mode.
void setup() {
  pinMode(MOM_PIN, INPUT_PULLDOWN);
  pinMode(LED_PIN, OUTPUT);
  Serial.begin(115200);
}

// Main loop, check the switch state
void loop() {

  int momState = digitalRead(MOM_PIN);

  // Only act if the state changed from the last run.
  if (previousMomState != momState) {
    Serial.printf("Arcade: %d\n", momState);
    previousMomState = momState;
    if (momState == 1) {
      Particle.publish("simulate-catastrophe", String(momState), 60, PRIVATE);\
      blink(125);
      blink(125);
      blink(125);
    }
  }

  // Might miss a state change if someone presses the switch and releases within the delay.
  // But also, seems excessive to not sleep between runs.
  // Realistically, who will press the button that fast?
  delay(200);
}

/*
 * Blinks the LED once with the given delay.
 */
void blink(unsigned long delayMs) {
  digitalWrite(LED_PIN, HIGH);
  delay(delayMs);
  digitalWrite(LED_PIN, LOW);
  delay(delayMs);
}
