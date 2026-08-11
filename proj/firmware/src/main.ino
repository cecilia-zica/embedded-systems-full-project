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

const int PINO_LED = 13; // active-high: HIGH = on, LOW = off

int bpmThreshold = 40; // initial value until the first backend response
bool alertEnabled = true;

// 50000 alone can't tell a finger from a surface reflecting light; it's just a
// cheap gate. The real check is dedoConfirmado() (signal variation).
const long IR_DEDO_PRESENTE = 50000;
// ~4s: a heartbeat lasts 0.6-1s, enough to catch 1-2 cycles before deciding.
const unsigned long DURACAO_CONFIRMACAO_MS = 4000;
const int INTERVALO_AMOSTRA_MS = 20;
// A finger varies the IR with the heartbeat; a table/box reflects steady light.
// This module's noise alone already reached 500-2900 (measured with no finger),
// hence 3000.
const long VARIACAO_MINIMA = 3000;
const unsigned long INTERVALO_ENVIO = 5000;
const unsigned long CONFIG_INTERVAL_MS = 10000; // how often to poll the backend config

const int MEIO_PERIODO_RAPIDO = 80; // fast blink = inconsistent reading
const int MEIO_PERIODO_LENTO = 250; // slow blink = normal reading
const int ALERTA_BLINK_ON = 100;    // 2 blinks + pause = alert
const int ALERTA_PAUSA = 400;

MAX30105 particleSensor;

long lastBeat = 0;
float bpmAtual = 0;
unsigned long ultimaConfigMs = 0;

// Noise also triggers checkForBeat but jumps around between beats
// (61, 107, 86, 130... measured with no finger); only accept after RATE_SIZE
// beats agree.
const byte RATE_SIZE = 4;
float taxasBpm[RATE_SIZE];
byte indiceTaxa = 0;
byte contagemTaxas = 0;

void ledOn() { digitalWrite(PINO_LED, HIGH); }
void ledOff() { digitalWrite(PINO_LED, LOW); }

// Cheap check used inside the read loop so it doesn't block; it can't tell a
// finger from a surface reflecting light.
bool dedoPresente() {
  return particleSensor.getIR() >= IR_DEDO_PRESENTE;
}

// GET /controle: refresh bpmThreshold/alertEnabled with what the app saved.
// Parsed by hand because there are only two fixed fields; a JSON library isn't
// worth the cost.
void buscarConfig() {
  if (WiFi.status() != WL_CONNECTED) return;

  HTTPClient http;
  String url = String("http://") + BACKEND_HOST + ":" + BACKEND_PORT + "/api/v1/controle";
  http.begin(url);
  http.addHeader("X-API-Key", API_KEY);

  int status = http.GET();
  if (status == 200) {
    String corpo = http.getString();
    int i = corpo.indexOf("\"bpm_threshold\":");
    if (i != -1) bpmThreshold = corpo.substring(i + 16).toInt();
    alertEnabled = corpo.indexOf("\"alert_enabled\":true") != -1;
    Serial.printf("CONFIG: bpmThreshold=%d alertEnabled=%d\n", bpmThreshold, alertEnabled);
  } else {
    Serial.printf("CONFIG: fetch failed (status=%d), keeping current values\n", status);
  }
  http.end();
}

// Only refetches after CONFIG_INTERVAL_MS; called from within the states so the
// sensor never blocks waiting on the network.
void atualizarConfigSePreciso() {
  if (millis() - ultimaConfigMs < CONFIG_INTERVAL_MS) return;
  ultimaConfigMs = millis();
  buscarConfig();
}

// The real confirmation (only in semDedo, needn't be fast): watch ~4s; a real
// finger varies the IR with the heartbeat, a table/box stays constant.
bool dedoConfirmado() {
  if (!dedoPresente()) return false;

  long minIR = 999999, maxIR = 0; // MAX30105 is an 18-bit ADC, IR never exceeds ~262143
  unsigned long fim = millis() + DURACAO_CONFIRMACAO_MS;
  while (millis() < fim) {
    if (!dedoPresente()) return false; // removed mid-way, cancel the confirmation

    long amostra = particleSensor.getIR();
    if (amostra < minIR) minIR = amostra;
    if (amostra > maxIR) maxIR = amostra;

    atualizarConfigSePreciso();
    delay(INTERVALO_AMOSTRA_MS);
  }
  return (maxIR - minIR) >= VARIACAO_MINIMA;
}

void piscarRitmo(int meioPeriodoMs, unsigned long duracaoMs) {
  unsigned long fim = millis() + duracaoMs;
  while (millis() < fim) {
    ledOn();
    delay(meioPeriodoMs);
    ledOff();
    delay(meioPeriodoMs);
  }
}

void piscarAlerta(unsigned long duracaoMs) {
  unsigned long fim = millis() + duracaoMs;
  while (millis() < fim) {
    for (int i = 0; i < 2; i++) {
      ledOn();
      delay(ALERTA_BLINK_ON);
      ledOff();
      delay(ALERTA_BLINK_ON);
    }
    delay(ALERTA_PAUSA);
  }
}

// Rough estimate, not the official clinical algorithm.
float estimarSpO2() {
  long redValue = particleSensor.getRed();
  long irValue = particleSensor.getIR();
  if (irValue == 0) return 0;

  float ratio = (float)redValue / (float)irValue;
  float spo2 = 110.0 - 25.0 * ratio;
  if (spo2 > 100) spo2 = 100;
  if (spo2 < 80) spo2 = 80;
  return spo2;
}

// Fixed rule for now (see TINYML.md); SpO2 is excluded because estimarSpO2()
// reads 85-87% even with a good signal, not reliable enough.
int classificar(float bpm, float spo2) {
  if (bpm < 15.0 || bpm > 220.0) return 2;          // error
  if (alertEnabled && bpm > bpmThreshold) return 1; // alert
  return 0;                                         // normal
}

// ---- states ----

// Cheap wait with no IR; once a finger appears, confirm for ~4s (dedoConfirmado)
// before accepting it.
void semDedo() {
  ledOff();
  while (!dedoConfirmado()) {
    atualizarConfigSePreciso();
    delay(50);
  }
  ledOn(); // finger confirmed: LED stays on while reading the heartbeat
}

// Requires RATE_SIZE consecutive beats within 30% of the mean, filtering noise
// that slipped past dedoConfirmado.
void lendoBatimento() {
  bpmAtual = 0;
  contagemTaxas = 0;

  while (bpmAtual == 0) {
    atualizarConfigSePreciso();
    if (!dedoPresente()) return; // finger removed: back to semDedo

    if (checkForBeat(particleSensor.getIR())) {
      long agora = millis();
      float bpmInstantaneo = 60000.0 / (agora - lastBeat);
      lastBeat = agora;

      if (bpmInstantaneo > 20 && bpmInstantaneo < 255) {
        taxasBpm[indiceTaxa++] = bpmInstantaneo;
        indiceTaxa %= RATE_SIZE;
        if (contagemTaxas < RATE_SIZE) contagemTaxas++;

        if (contagemTaxas == RATE_SIZE) {
          float soma = 0, minTaxa = 999, maxTaxa = 0;
          for (byte i = 0; i < RATE_SIZE; i++) {
            soma += taxasBpm[i];
            if (taxasBpm[i] < minTaxa) minTaxa = taxasBpm[i];
            if (taxasBpm[i] > maxTaxa) maxTaxa = taxasBpm[i];
          }
          float media = soma / RATE_SIZE;
          if ((maxTaxa - minTaxa) <= media * 0.3) {
            bpmAtual = media;
          } else {
            contagemTaxas = 0; // inconsistent, discard and restart the series
          }
        }
      }
    }
  }
}

// POST to the backend; with no WiFi it only warns on serial (the reading was
// already printed above).
void enviarLog(float bpm, float spo2, int classe) {
  if (WiFi.status() != WL_CONNECTED) {
    Serial.println("-> backend: no WiFi, not sent (serial only).");
    return;
  }

  HTTPClient http;
  String url = String("http://") + BACKEND_HOST + ":" + BACKEND_PORT + "/api/v1/logging";
  http.begin(url);
  http.addHeader("Content-Type", "application/json");
  http.addHeader("X-API-Key", API_KEY);

  char corpo[128];
  snprintf(corpo, sizeof(corpo),
           "{\"bpm\":%.1f,\"spo2\":%.1f,\"class\":%d}",
           bpm, spo2, classe);

  int status = http.POST(corpo);
  Serial.printf("-> backend: status=%d\n", status);
  http.end();
}

void classificarEEnviar() {
  float spo2 = estimarSpO2();
  int classe = classificar(bpmAtual, spo2);

  Serial.printf("READING,%.1f,%.1f,%d\n", bpmAtual, spo2, classe);
  enviarLog(bpmAtual, spo2, classe);

  switch (classe) {
    case 2: piscarRitmo(MEIO_PERIODO_RAPIDO, INTERVALO_ENVIO); break;
    case 1: piscarAlerta(INTERVALO_ENVIO); break;
    default: piscarRitmo(MEIO_PERIODO_LENTO, INTERVALO_ENVIO);
  }
}

// Connects with a timeout; on failure it runs offline (sensor + serial keep going).
void conectarWiFi() {
  WiFi.mode(WIFI_STA);
  WiFi.setTxPower(WIFI_POWER_8_5dBm); // lowers the radio current peak (see note on top)
  WiFi.begin(WIFI_SSID, WIFI_SENHA);

  Serial.printf("Connecting to WiFi '%s'", WIFI_SSID);
  unsigned long inicio = millis();
  while (WiFi.status() != WL_CONNECTED && millis() - inicio < WIFI_TIMEOUT_MS) {
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
  pinMode(PINO_LED, OUTPUT);
  ledOff();

  conectarWiFi();
  buscarConfig(); // initial config, before the first read cycle
  ultimaConfigMs = millis();

  Wire.begin();
  Wire.setTimeOut(200); // avoids blocking forever if a wire comes loose

  if (!particleSensor.begin(Wire, I2C_SPEED_STANDARD)) {
    Serial.println("ERROR: MAX30102 not found!");
    while (1) { // blink fast forever = broken hardware
      ledOn();
      delay(MEIO_PERIODO_RAPIDO);
      ledOff();
      delay(MEIO_PERIODO_RAPIDO);
    }
  }
  particleSensor.setup();
  particleSensor.setPulseAmplitudeRed(0x1F);
  particleSensor.setPulseAmplitudeIR(0x1F);

  Serial.println("Ready. Place your finger to measure.");
}

// Old bug: every reading returned to the full semDedo(), which reconfirms ~4s of
// variation; a still finger might never cross the threshold and the LED stayed
// off, looking stuck. Now it confirms once and keeps reading while dedoPresente().
void loop() {
  semDedo();
  while (dedoPresente()) {
    lendoBatimento();
    if (bpmAtual == 0) break; // finger left mid-reading
    classificarEEnviar();
  }
}

// Missing: real TinyML (see TINYML.md), clinical-grade SpO2.
