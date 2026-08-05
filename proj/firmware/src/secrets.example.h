#pragma once

// Template de credenciais. Copie este arquivo para "secrets.h" (no mesmo
// diretório) e preencha com os valores reais. secrets.h é ignorado pelo git
// (ver proj/firmware/.gitignore), então senha de WiFi e chave de API nunca vão
// parar no histórico do repositório.
//
//   cp src/secrets.example.h src/secrets.h

// Rede WiFi que o ESP32 usa (na apresentação, o hotspot do celular).
const char *WIFI_SSID  = "SUA_REDE_WIFI";
const char *WIFI_SENHA = "SUA_SENHA_WIFI";

// IP e porta do backend Go na rede local. Ache o IP com `ipconfig getifaddr en0`.
const char *BACKEND_HOST = "192.168.0.10";
const int   BACKEND_PORT = 8080;

// Tem que bater com a variável de ambiente API_KEY do backend
// (ou o defaultAPIKey em proj/backend/middleware.go).
const char *API_KEY = "TROQUE_ESTA_CHAVE";
