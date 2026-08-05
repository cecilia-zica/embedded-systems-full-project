//toda lógica HTTP aqui — telas não sabem nada de HTTP, só chamam esses métodos
//equivalente ao api.js/axios.js do RN; http == axios/fetch, tudo devolve Future<T> (Promise<T>)

import 'dart:convert';
import 'package:http/http.dart' as http;

class ApiService {
  // baseUrl e apiKey vêm de --dart-define no build, não ficam chumbados no
  // código: o mesmo binário roda em qualquer rede/ambiente só trocando o build.
  //   flutter run --dart-define=API_BASE_URL=http://SEU_IP:8080 --dart-define=API_KEY=...
  // Os defaults servem só pro dev local (hotspot do Mac).
  static const String baseUrl = String.fromEnvironment(
    'API_BASE_URL',
    defaultValue: 'http://172.20.10.2:8080',
  );

  //tem que bater com a API_KEY do backend (defaultAPIKey em middleware.go no dev)
  static const String apiKey = String.fromEnvironment(
    'API_KEY',
    defaultValue: 'zica123',
  );

  // todas as rotas exigem X-API-Key agora (leitura e escrita)
  static const Map<String, String> _headers = {
    'Content-Type': 'application/json',
    'X-API-Key': apiKey,
  };

  // GET /api/v1/logging — busca o histórico de leituras
  static Future<List<dynamic>> getLogs() async {
    final response = await http.get(
      Uri.parse('$baseUrl/api/v1/logging'),
      headers: _headers,
    );
    if (response.statusCode != 200) {
      throw Exception('erro ao buscar logs: ${response.statusCode}');
    }
    //json.decode: JSON -> List/Map do Dart, mesma ideia do Decode do Go
    return json.decode(response.body) as List<dynamic>;
  }

  // GET /api/v1/controle — busca a config atual (bpm_threshold, alert_enabled)
  static Future<Map<String, dynamic>> getConfig() async {
    final response = await http.get(
      Uri.parse('$baseUrl/api/v1/controle'),
      headers: _headers,
    );
    if (response.statusCode != 200) {
      throw Exception('erro ao buscar config: ${response.statusCode}');
    }
    return json.decode(response.body) as Map<String, dynamic>;
  }

  // DELETE /api/v1/logging — apaga todo o histórico de leituras
  static Future<void> deleteLogs() async {
    final response = await http.delete(
      Uri.parse('$baseUrl/api/v1/logging'),
      headers: _headers,
    );
    if (response.statusCode != 200) {
      throw Exception('erro ao apagar logs: ${response.statusCode}');
    }
  }

  // POST /api/v1/controle — salva um novo threshold/alerta
  static Future<void> postConfig(int threshold, bool alertEnabled) async {
    final response = await http.post(
      Uri.parse('$baseUrl/api/v1/controle'),
      headers: _headers,
      //json.encode: Map do Dart -> JSON, mesma ideia do Encode do Go
      body: json.encode({
        'bpm_threshold': threshold,
        'alert_enabled': alertEnabled,
      }),
    );
    if (response.statusCode != 200) {
      throw Exception('erro ao salvar config: ${response.statusCode}');
    }
  }
}
