#pragma once

// Credentials template. Copy this file to "secrets.h" (same directory) and fill
// in the real values. secrets.h is git-ignored (see proj/firmware/.gitignore),
// so the WiFi password and API key never reach the repository history.
//
//   cp src/secrets.example.h src/secrets.h

// WiFi network the ESP32 joins (a phone hotspot during the demo).
const char *WIFI_SSID  = "YOUR_WIFI_SSID";
const char *WIFI_SENHA = "YOUR_WIFI_PASSWORD";

// Backend IP and port on the local network. Find the IP with `ipconfig getifaddr en0`.
const char *BACKEND_HOST = "192.168.0.10";
const int   BACKEND_PORT = 8080;

// Must match the backend's API_KEY environment variable
// (or defaultAPIKey in proj/backend/middleware.go).
const char *API_KEY = "CHANGE_THIS_KEY";
