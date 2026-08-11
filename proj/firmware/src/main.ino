// ESP32 + MAX30102 + WiFi: reads BPM/SpO2, classifies it and POSTs to the backend.
// bpmThreshold/alertEnabled are refreshed from the backend every
// CONFIG_INTERVAL_MS (GET /controle).
// RFID was removed: out of scope and the module never communicated reliably.
// WiFi: a reduced setTxPower avoids stalling the I2C bus; if it stalls again the
// cause is power, not code.
// Demo network = phone hotspot; the IP changes each time (ipconfig getifaddr en0).

#include <WiFi.h>
#include <HTTPClient.h>
#include <Wire.h>
#include "MAX30105.h"
#include "heartRate.h"

// Credentials (WiFi, backend host, API key) come from secrets.h, a local file
// kept out of git. Copy secrets.example.h to secrets.h and fill it in.
#include "secrets.h"

const unsigned long WIFI_TIMEOUT_MS = 15000; // give up after this and run offline

const int LED_PIN = 13; // active-high: HIGH = on, LOW = off

int bpmThreshold = 40; // initial value until the first backend response
bool alertEnabled = true;

// 50000 alone can't tell a finger from a surface reflecting light; it's just a
// cheap gate. The real check is fingerConfirmed() (signal variation).
const long IR_FINGER_PRESENT = 50000;
// ~4s: a heartbeat lasts 0.6-1s, enough to catch 1-2 cycles before deciding.
const unsigned long CONFIRM_DURATION_MS = 4000;
const int SAMPLE_INTERVAL_MS = 20;
// A finger varies the IR with the heartbeat; a table/box reflects steady light.
// This module's noise alone already reached 500-2900 (measured with no finger),
// hence 3000.
const long MIN_VARIATION = 3000;
const unsigned long SEND_INTERVAL = 5000;
const unsigned long CONFIG_INTERVAL_MS = 10000; // how often to poll the backend config

const int HALF_PERIOD_FAST = 80; // fast blink = inconsistent reading
const int HALF_PERIOD_SLOW = 250; // slow blink = normal reading
const int ALERT_BLINK_ON = 100;    // 2 blinks + pause = alert
const int ALERT_PAUSE = 400;

MAX30105 particleSensor;

long lastBeat = 0;
float currentBpm = 0;
unsigned long lastConfigMs = 0;

// Noise also triggers checkForBeat but jumps around between beats
// (61, 107, 86, 130... measured with no finger); only accept after RATE_SIZE
// beats agree.
const byte RATE_SIZE = 4;
float bpmRates[RATE_SIZE];
byte rateIndex = 0;
byte rateCount = 0;

void ledOn() { digitalWrite(LED_PIN, HIGH); }
void ledOff() { digitalWrite(LED_PIN, LOW); }

// Cheap check used inside the read loop so it doesn't block; it can't tell a
// finger from a surface reflecting light.
bool fingerPresent() {
  return particleSensor.getIR() >= IR_FINGER_PRESENT;
}

// GET /controle: refresh bpmThreshold/alertEnabled with what the app saved.
// Parsed by hand because there are only two fixed fields; a JSON library isn't
// worth the cost.
void fetchConfig() {
  if (WiFi.status() != WL_CONNECTED) return;

  HTTPClient http;
  String url = String("http://") + BACKEND_HOST + ":" + BACKEND_PORT + "/api/v1/controle";
  http.begin(url);
  http.addHeader("X-API-Key", API_KEY);

  int status = http.GET();
  if (status == 200) {
    String body = http.getString();
    int i = body.indexOf("\"bpm_threshold\":");
    if (i != -1) bpmThreshold = body.substring(i + 16).toInt();
    alertEnabled = body.indexOf("\"alert_enabled\":true") != -1;
    Serial.printf("CONFIG: bpmThreshold=%d alertEnabled=%d\n", bpmThreshold, alertEnabled);
  } else {
    Serial.printf("CONFIG: fetch failed (status=%d), keeping current values\n", status);
  }
  http.end();
}

// Only refetches after CONFIG_INTERVAL_MS; called from within the states so the
// sensor never blocks waiting on the network.
void refreshConfigIfDue() {
  if (millis() - lastConfigMs < CONFIG_INTERVAL_MS) return;
  lastConfigMs = millis();
  fetchConfig();
}

// The real confirmation (only in waitForFinger, needn't be fast): watch ~4s; a real
// finger varies the IR with the heartbeat, a table/box stays constant.
bool fingerConfirmed() {
  if (!fingerPresent()) return false;

  long minIR = 999999, maxIR = 0; // MAX30105 is an 18-bit ADC, IR never exceeds ~262143
  unsigned long end = millis() + CONFIRM_DURATION_MS;
  while (millis() < end) {
    if (!fingerPresent()) return false; // removed mid-way, cancel the confirmation

    long sample = particleSensor.getIR();
    if (sample < minIR) minIR = sample;
    if (sample > maxIR) maxIR = sample;

    refreshConfigIfDue();
    delay(SAMPLE_INTERVAL_MS);
  }
  return (maxIR - minIR) >= MIN_VARIATION;
}

void blinkRhythm(int halfPeriodMs, unsigned long durationMs) {
  unsigned long end = millis() + durationMs;
  while (millis() < end) {
    ledOn();
    delay(halfPeriodMs);
    ledOff();
    delay(halfPeriodMs);
  }
}

void blinkAlert(unsigned long durationMs) {
  unsigned long end = millis() + durationMs;
  while (millis() < end) {
    for (int i = 0; i < 2; i++) {
      ledOn();
      delay(ALERT_BLINK_ON);
      ledOff();
      delay(ALERT_BLINK_ON);
    }
    delay(ALERT_PAUSE);
  }
}

// Rough estimate, not the official clinical algorithm.
float estimateSpO2() {
  long redValue = particleSensor.getRed();
  long irValue = particleSensor.getIR();
  if (irValue == 0) return 0;

  float ratio = (float)redValue / (float)irValue;
  float spo2 = 110.0 - 25.0 * ratio;
  if (spo2 > 100) spo2 = 100;
  if (spo2 < 80) spo2 = 80;
  return spo2;
}

// Fixed rule for now (see TINYML.md); SpO2 is excluded because estimateSpO2()
// reads 85-87% even with a good signal, not reliable enough.
int classify(float bpm, float spo2) {
  if (bpm < 15.0 || bpm > 220.0) return 2;          // error
  if (alertEnabled && bpm > bpmThreshold) return 1; // alert
  return 0;                                         // normal
}

// ---- states ----

// Cheap wait with no IR; once a finger appears, confirm for ~4s (fingerConfirmed)
// before accepting it.
void waitForFinger() {
  ledOff();
  while (!fingerConfirmed()) {
    refreshConfigIfDue();
    delay(50);
  }
  ledOn(); // finger confirmed: LED stays on while reading the heartbeat
}

// Requires RATE_SIZE consecutive beats within 30% of the mean, filtering noise
// that slipped past fingerConfirmed.
void readHeartbeat() {
  currentBpm = 0;
  rateCount = 0;

  while (currentBpm == 0) {
    refreshConfigIfDue();
    if (!fingerPresent()) return; // finger removed: back to waitForFinger

    if (checkForBeat(particleSensor.getIR())) {
      long now = millis();
      float instantBpm = 60000.0 / (now - lastBeat);
      lastBeat = now;

      if (instantBpm > 20 && instantBpm < 255) {
        bpmRates[rateIndex++] = instantBpm;
        rateIndex %= RATE_SIZE;
        if (rateCount < RATE_SIZE) rateCount++;

        if (rateCount == RATE_SIZE) {
          float sum = 0, minRate = 999, maxRate = 0;
          for (byte i = 0; i < RATE_SIZE; i++) {
            sum += bpmRates[i];
            if (bpmRates[i] < minRate) minRate = bpmRates[i];
            if (bpmRates[i] > maxRate) maxRate = bpmRates[i];
          }
          float mean = sum / RATE_SIZE;
          if ((maxRate - minRate) <= mean * 0.3) {
            currentBpm = mean;
          } else {
            rateCount = 0; // inconsistent, discard and restart the series
          }
        }
      }
    }
  }
}

// POST to the backend; with no WiFi it only warns on serial (the reading was
// already printed above).
void sendLog(float bpm, float spo2, int cls) {
  if (WiFi.status() != WL_CONNECTED) {
    Serial.println("-> backend: no WiFi, not sent (serial only).");
    return;
  }

  HTTPClient http;
  String url = String("http://") + BACKEND_HOST + ":" + BACKEND_PORT + "/api/v1/logging";
  http.begin(url);
  http.addHeader("Content-Type", "application/json");
  http.addHeader("X-API-Key", API_KEY);

  char body[128];
  snprintf(body, sizeof(body),
           "{\"bpm\":%.1f,\"spo2\":%.1f,\"class\":%d}",
           bpm, spo2, cls);

  int status = http.POST(body);
  Serial.printf("-> backend: status=%d\n", status);
  http.end();
}

void classifyAndSend() {
  float spo2 = estimateSpO2();
  int cls = classify(currentBpm, spo2);

  Serial.printf("READING,%.1f,%.1f,%d\n", currentBpm, spo2, cls);
  sendLog(currentBpm, spo2, cls);

  switch (cls) {
    case 2: blinkRhythm(HALF_PERIOD_FAST, SEND_INTERVAL); break;
    case 1: blinkAlert(SEND_INTERVAL); break;
    default: blinkRhythm(HALF_PERIOD_SLOW, SEND_INTERVAL);
  }
}

// Connects with a timeout; on failure it runs offline (sensor + serial keep going).
void connectWiFi() {
  WiFi.mode(WIFI_STA);
  WiFi.setTxPower(WIFI_POWER_8_5dBm); // lowers the radio current peak (see note on top)
  WiFi.begin(WIFI_SSID, WIFI_PASSWORD);

  Serial.printf("Connecting to WiFi '%s'", WIFI_SSID);
  unsigned long start = millis();
  while (WiFi.status() != WL_CONNECTED && millis() - start < WIFI_TIMEOUT_MS) {
    delay(300);
    Serial.print(".");
  }
  Serial.println();

  if (WiFi.status() == WL_CONNECTED) {
    Serial.print("WiFi connected! IP: ");
    Serial.println(WiFi.localIP());
  } else {
    Serial.println("WARNING: WiFi did not connect (timeout) - running offline.");
  }
}

void setup() {
  Serial.begin(115200);
  pinMode(LED_PIN, OUTPUT);
  ledOff();

  connectWiFi();
  fetchConfig(); // initial config, before the first read cycle
  lastConfigMs = millis();

  Wire.begin();
  Wire.setTimeOut(200); // avoids blocking forever if a wire comes loose

  if (!particleSensor.begin(Wire, I2C_SPEED_STANDARD)) {
    Serial.println("ERROR: MAX30102 not found!");
    while (1) { // blink fast forever = broken hardware
      ledOn();
      delay(HALF_PERIOD_FAST);
      ledOff();
      delay(HALF_PERIOD_FAST);
    }
  }
  particleSensor.setup();
  particleSensor.setPulseAmplitudeRed(0x1F);
  particleSensor.setPulseAmplitudeIR(0x1F);

  Serial.println("Ready. Place your finger to measure.");
}

// Old bug: every reading returned to the full waitForFinger(), which reconfirms ~4s of
// variation; a still finger might never cross the threshold and the LED stayed
// off, looking stuck. Now it confirms once and keeps reading while fingerPresent().
void loop() {
  waitForFinger();
  while (fingerPresent()) {
    readHeartbeat();
    if (currentBpm == 0) break; // finger left mid-reading
    classifyAndSend();
  }
}

// Missing: real TinyML (see TINYML.md), clinical-grade SpO2.
