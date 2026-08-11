import 'dart:convert';

import 'package:http/http.dart' as http;

/// HTTP client for the monitoring backend.
///
/// Centralizes every API call so widgets depend only on these methods. All
/// routes require the `X-API-Key` header.
class ApiService {
  /// Backend base URL, injected at build time via `--dart-define=API_BASE_URL`.
  ///
  /// The default targets local development and can be overridden per build.
  static const String baseUrl = String.fromEnvironment(
    'API_BASE_URL',
    defaultValue: 'http://172.20.10.2:8080',
  );

  /// API key sent on every request, injected via `--dart-define=API_KEY`.
  ///
  /// Must match the backend's `API_KEY`; the default is for local development.
  static const String apiKey = String.fromEnvironment(
    'API_KEY',
    defaultValue: 'zica123',
  );

  static const Map<String, String> _headers = {
    'Content-Type': 'application/json',
    'X-API-Key': apiKey,
  };

  /// Fetches the reading history from `GET /api/v1/logging`.
  static Future<List<dynamic>> getLogs() async {
    final response = await http.get(
      Uri.parse('$baseUrl/api/v1/logging'),
      headers: _headers,
    );
    if (response.statusCode != 200) {
      throw Exception('failed to fetch logs: ${response.statusCode}');
    }
    return json.decode(response.body) as List<dynamic>;
  }

  /// Fetches the current configuration from `GET /api/v1/controle`.
  static Future<Map<String, dynamic>> getConfig() async {
    final response = await http.get(
      Uri.parse('$baseUrl/api/v1/controle'),
      headers: _headers,
    );
    if (response.statusCode != 200) {
      throw Exception('failed to fetch config: ${response.statusCode}');
    }
    return json.decode(response.body) as Map<String, dynamic>;
  }

  /// Deletes the entire reading history via `DELETE /api/v1/logging`.
  static Future<void> deleteLogs() async {
    final response = await http.delete(
      Uri.parse('$baseUrl/api/v1/logging'),
      headers: _headers,
    );
    if (response.statusCode != 200) {
      throw Exception('failed to delete logs: ${response.statusCode}');
    }
  }

  /// Saves a new alert threshold and toggle via `POST /api/v1/controle`.
  static Future<void> postConfig(int threshold, bool alertEnabled) async {
    final response = await http.post(
      Uri.parse('$baseUrl/api/v1/controle'),
      headers: _headers,
      body: json.encode({
        'bpm_threshold': threshold,
        'alert_enabled': alertEnabled,
      }),
    );
    if (response.statusCode != 200) {
      throw Exception('failed to save config: ${response.statusCode}');
    }
  }
}
