//ESP32+MAX30102+WiFi: lê BPM/SpO2, classifica e manda pro backend via HTTP
//bpmThreshold/alertEnabled vêm do backend a cada CONFIG_INTERVAL_MS (GET /controle)
//RFID removido 08/07: fora da rubrica, módulo nunca comunicou direito (ver ARDUINO.md)
//WiFi: setTxPower reduzido evita travar o I2C; travou de novo = alimentação, não código
//rede da apresentação = hotspot do iPhone; IP muda toda vez, pegar com `ipconfig getifaddr en0`

#include <WiFi.h>
#include <HTTPClient.h>
#include <Wire.h>
#include "MAX30105.h"
#include "heartRate.h"

// Credenciais (WiFi, host do backend, API key) vêm de secrets.h — arquivo
// local, fora do git. Copie secrets.example.h para secrets.h e preencha.
#include "secrets.h"

const unsigned long WIFI_TIMEOUT_MS = 15000; // desiste depois disso e segue sem WiFi

const int PINO_LED = 13; // ativo-ALTO: HIGH = aceso, LOW = apagado (confirmado testando)

int bpmThreshold = 40; // valor inicial até a 1a resposta do backend (ver buscarConfig)
bool alertEnabled = true;

//50000 sozinho não distingue dedo de mesa refletindo luz; é só filtro barato,
//confirmação de verdade é dedoConfirmado() (variação do sinal)
const long IR_DEDO_PRESENTE = 50000;
//~4s: batimento dura 0,6-1s, dá tempo de pegar 1-2 ciclos antes de decidir
const unsigned long DURACAO_CONFIRMACAO_MS = 4000;
const int INTERVALO_AMOSTRA_MS = 20;
//dedo varia com o batimento, mesa/caixa reflete luz constante; ruído desse
//módulo já bateu 500-2900 sozinho (medido sem dedo), por isso 3000
const long VARIACAO_MINIMA = 3000;
const unsigned long INTERVALO_ENVIO = 5000;
const unsigned long CONFIG_INTERVAL_MS = 10000; // frequência de busca do bpm_threshold/alert_enabled no backend (mais rápido pra dar pra ver mudando ao vivo na apresentação)

const int MEIO_PERIODO_RAPIDO = 80; // pisca rápido = leitura inconsistente
const int MEIO_PERIODO_LENTO = 250; // pisca devagar = leitura normal
const int ALERTA_BLINK_ON = 100;    // 2 piscadas + pausa = alerta
const int ALERTA_PAUSA = 400;

MAX30105 particleSensor;

long lastBeat = 0;
float bpmAtual = 0;
unsigned long ultimaConfigMs = 0;

//ruído também dispara checkForBeat mas pula de valor a cada batida
//(61, 107, 86, 130... testado sem dedo); só aceita após RATE_SIZE batidas concordando
const byte RATE_SIZE = 4;
float taxasBpm[RATE_SIZE];
byte indiceTaxa = 0;
byte contagemTaxas = 0;

void ledOn() { digitalWrite(PINO_LED, HIGH); }
void ledOff() { digitalWrite(PINO_LED, LOW); }

//checagem rápida/barata, usada dentro da leitura pra não travar;
//não distingue dedo de superfície refletindo luz
bool dedoPresente() {
  return particleSensor.getIR() >= IR_DEDO_PRESENTE;
}

//GET /controle: atualiza bpmThreshold/alertEnabled com o que o app salvou.
//Parse manual porque são só 2 campos fixos, não compensa trazer lib de JSON
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
    Serial.printf("CONFIG: falha ao buscar (status=%d), mantendo valores atuais\n", status);
  }
  http.end();
}

//só busca de novo após CONFIG_INTERVAL_MS; chamada de dentro dos estados pra não travar o sensor esperando rede
void atualizarConfigSePreciso() {
  if (millis() - ultimaConfigMs < CONFIG_INTERVAL_MS) return;
  ultimaConfigMs = millis();
  buscarConfig();
}

//confirmação de verdade (só em semDedo, não precisa ser rápida): observa ~4s,
//dedo de verdade varia o IR com o batimento, mesa/caixa fica constante
bool dedoConfirmado() {
  if (!dedoPresente()) return false;

  long minIR = 999999, maxIR = 0; // MAX30105 é ADC de 18 bits, IR nunca passa de ~262143
  unsigned long fim = millis() + DURACAO_CONFIRMACAO_MS;
  while (millis() < fim) {
    if (!dedoPresente()) return false; // tirou no meio, cancela a confirmação

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

// estimativa simples (não é o algoritmo clínico oficial)
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

//regra fixa por enquanto (ver TINYML.md); SpO2 não entra pq estimarSpO2() sai
//85-87% mesmo com leitura boa, não é confiável o bastante
int classificar(float bpm, float spo2) {
  if (bpm < 15.0 || bpm > 220.0) return 2;          // erro
  if (alertEnabled && bpm > bpmThreshold) return 1; // alerta
  return 0;                                         // normal
}

// ---- estados ----

//espera barata sem IR; ao aparecer, confirma por ~4s (dedoConfirmado) antes de aceitar
void semDedo() {
  ledOff();
  while (!dedoConfirmado()) {
    atualizarConfigSePreciso();
    delay(50);
  }
  ledOn(); // dedo confirmado: LED aceso fixo enquanto lê o batimento
}

//exige RATE_SIZE batidas seguidas dentro de 30% da média — filtra ruído que passou do dedoConfirmado
void lendoBatimento() {
  bpmAtual = 0;
  contagemTaxas = 0;

  while (bpmAtual == 0) {
    atualizarConfigSePreciso();
    if (!dedoPresente()) return; // tirou o dedo: volta pro semDedo

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
            contagemTaxas = 0; // inconsistente, descarta e recomeça a série
          }
        }
      }
    }
  }
}

//POST pro backend; sem WiFi só avisa na serial, a leitura já foi impressa antes
void enviarLog(float bpm, float spo2, int classe) {
  if (WiFi.status() != WL_CONNECTED) {
    Serial.println("-> backend: sem WiFi, não enviei (só ficou na serial).");
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

  Serial.printf("LEITURA,%.1f,%.1f,%d\n", bpmAtual, spo2, classe);
  enviarLog(bpmAtual, spo2, classe);

  switch (classe) {
    case 2: piscarRitmo(MEIO_PERIODO_RAPIDO, INTERVALO_ENVIO); break;
    case 1: piscarAlerta(INTERVALO_ENVIO); break;
    default: piscarRitmo(MEIO_PERIODO_LENTO, INTERVALO_ENVIO);
  }
}

// conecta com timeout; sem sucesso, segue sem WiFi (sensor+serial continuam)
void conectarWiFi() {
  WiFi.mode(WIFI_STA);
  WiFi.setTxPower(WIFI_POWER_8_5dBm); // reduz o pico de corrente do rádio (ver nota no topo)
  WiFi.begin(WIFI_SSID, WIFI_SENHA);

  Serial.printf("Conectando no WiFi '%s'", WIFI_SSID);
  unsigned long inicio = millis();
  while (WiFi.status() != WL_CONNECTED && millis() - inicio < WIFI_TIMEOUT_MS) {
    delay(300);
    Serial.print(".");
  }
  Serial.println();

  if (WiFi.status() == WL_CONNECTED) {
    Serial.print("WiFi conectado! IP: ");
    Serial.println(WiFi.localIP());
  } else {
    Serial.println("AVISO: WiFi não conectou (timeout) — seguindo sem WiFi.");
  }
}

void setup() {
  Serial.begin(115200);
  pinMode(PINO_LED, OUTPUT);
  ledOff();

  conectarWiFi();
  buscarConfig(); // config inicial, antes do primeiro ciclo de leitura
  ultimaConfigMs = millis();

  Wire.begin();
  Wire.setTimeOut(200); // evita travar pra sempre se um fio ficar solto

  if (!particleSensor.begin(Wire, I2C_SPEED_STANDARD)) {
    Serial.println("ERRO: MAX30102 não encontrado!");
    while (1) { // pisca rápido pra sempre = hardware quebrado
      ledOn();
      delay(MEIO_PERIODO_RAPIDO);
      ledOff();
      delay(MEIO_PERIODO_RAPIDO);
    }
  }
  particleSensor.setup();
  particleSensor.setPulseAmplitudeRed(0x1F);
  particleSensor.setPulseAmplitudeIR(0x1F);

  Serial.println("Pronto. Encoste o dedo pra medir.");
}

//bug antigo: cada leitura voltava pro semDedo() inteiro, que reconfirma ~4s de
//variação; dedo parado podia nunca rebater o limiar e o LED ficava apagado
//parecendo travado. agora confirma uma vez e segue lendo enquanto dedoPresente()
void loop() {
  semDedo();
  while (dedoPresente()) {
    lendoBatimento();
    if (bpmAtual == 0) break; // dedo saiu no meio da leitura
    classificarEEnviar();
  }
}

// falta: TinyML de verdade (ver TINYML.md), SpO2 de precisão clínica
