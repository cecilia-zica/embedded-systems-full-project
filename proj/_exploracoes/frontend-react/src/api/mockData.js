// Dados de demonstração — usados quando o fetch pro backend falha (ver
// hooks/useLogs.js e hooks/useConfig.js). Existe pra você conseguir ver e
// mexer no app mesmo sem o backend Go rodando ou sem o ESP32 conectado —
// útil testando em outra rede, sem o hardware por perto, etc.
//
// Timestamps relativos a "agora" (Date.now() - X minutos) em vez de datas
// fixas: assim os logs sempre parecem recentes, não importa quando você
// abrir o app.
function minutesAgo(minutes) {
  return new Date(Date.now() - minutes * 60_000).toISOString();
}

export const mockLogs = [
  { id: 8, bpm: 76, spo2: 98, class: 0, user_id: 'A1B2C3D4', created_at: minutesAgo(2) },
  { id: 7, bpm: 132, spo2: 95, class: 1, user_id: 'A1B2C3D4', created_at: minutesAgo(9) },
  { id: 6, bpm: 81, spo2: 97, class: 0, user_id: 'F09E1122', created_at: minutesAgo(18) },
  { id: 5, bpm: 0, spo2: 0, class: 2, user_id: '7C3A9B00', created_at: minutesAgo(25) },
  { id: 4, bpm: 74, spo2: 98, class: 0, user_id: 'A1B2C3D4', created_at: minutesAgo(40) },
  { id: 3, bpm: 145, spo2: 93, class: 1, user_id: 'D41D8CD9', created_at: minutesAgo(55) },
  { id: 2, bpm: 79, spo2: 99, class: 0, user_id: 'F09E1122', created_at: minutesAgo(70) },
  { id: 1, bpm: 88, spo2: 97, class: 0, user_id: 'A1B2C3D4', created_at: minutesAgo(95) },
];

export const mockConfig = { bpm_threshold: 120, alert_enabled: true };
