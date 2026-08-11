// Toda a lógica HTTP mora aqui — as telas (e agora os hooks em src/hooks/)
// não sabem nada sobre fetch/URLs, só chamam essas funções. Mesmo papel do
// services/api_service.dart no app Flutter (proj/app).

import { API_BASE_URL, API_KEY } from '../config';

const HEADERS = {
  'Content-Type': 'application/json',
  'X-API-Key': API_KEY,
};

// quanto tempo esperar antes de desistir de uma requisição. Sem isso, um IP
// inalcançável (backend desligado, rede errada) deixa o fetch pendurado por
// mais de 1 minuto antes do navegador desistir sozinho — tempo de sobra pra
// simplesmente tentar de novo mais tarde, mas ruim pra UI, que fica com
// spinner infinito em vez de cair pro modo demonstração rápido.
const REQUEST_TIMEOUT_MS = 4000;

// função interna (não exportada) que centraliza o "se não deu 2xx, joga erro"
async function request(path, options = {}) {
  // AbortController é a forma padrão do navegador de cancelar um fetch em
  // andamento; combinado com setTimeout, vira um "timeout" pro fetch (que
  // não tem essa opção pronta, ao contrário do timeout do http do Dart).
  const controller = new AbortController();
  const timeoutId = setTimeout(() => controller.abort(), REQUEST_TIMEOUT_MS);

  try {
    const res = await fetch(`${API_BASE_URL}${path}`, { ...options, headers: HEADERS, signal: controller.signal });
    if (!res.ok) {
      throw new Error(`erro em ${path}: ${res.status}`);
    }
    return res;
  } finally {
    clearTimeout(timeoutId);
  }
}

// GET /api/v1/logging — busca o histórico de leituras
export async function getLogs() {
  const res = await request('/api/v1/logging');
  return res.json();
}

// GET /api/v1/config — busca a configuração atual
export async function getConfig() {
  const res = await request('/api/v1/config');
  return res.json();
}

// POST /api/v1/config — salva um novo threshold/alerta
export async function postConfig(threshold, alertEnabled) {
  await request('/api/v1/config', {
    method: 'POST',
    body: JSON.stringify({ bpm_threshold: threshold, alert_enabled: alertEnabled }),
  });
}

// DELETE /api/v1/logging — apaga todo o histórico de leituras
export async function deleteLogs() {
  await request('/api/v1/logging', { method: 'DELETE' });
}
