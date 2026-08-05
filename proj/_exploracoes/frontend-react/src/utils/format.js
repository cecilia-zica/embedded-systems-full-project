// Helpers puros (sem estado, sem JSX). Diferença pra versão anterior: antes
// classToColor devolvia um hex cru (#2e7d32); agora classToStatus devolve
// uma PALAVRA semântica ('normal' | 'alert' | 'error'). Quem decide a cor de
// verdade é o CSS (ver theme/status.css) — assim o mesmo status pode ter
// tons diferentes no tema claro e no escuro sem essa função saber ou se
// importar com isso. É separar "o que o dado significa" de "como ele é
// desenhado", um princípio central de clean code.

// 0 = Repouso/Normal, 1 = Ativo/Alerta, 2 = Erro — mesma classificação do
// TinyML rodando no ESP32 (proj/firmware/src/classificacao.h)
export function classToStatus(c) {
  if (c === 0) return 'normal';
  if (c === 1) return 'alert';
  return 'error';
}

// O backend salva os timestamps em UTC (com "Z" no final). `new Date(iso)`
// já entende esse "Z" e os getters (getHours, getDate, ...) já devolvem o
// horário LOCAL do navegador automaticamente — mesmo resultado do
// `DateTime.parse(iso).toLocal()` do Dart, sem precisar converter fuso na mão.
export function formatTimestamp(iso) {
  const dt = new Date(iso);
  const two = (n) => String(n).padStart(2, '0');
  return `${two(dt.getDate())}/${two(dt.getMonth() + 1)} ${two(dt.getHours())}:${two(dt.getMinutes())}:${two(dt.getSeconds())}`;
}
