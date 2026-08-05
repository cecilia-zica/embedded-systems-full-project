// Componente "burro" (sem estado, sem hooks) — sempre desenha a mesma
// coisa. Ícone da lâmpada e os "dots" de status vêm da lib lucide-react em
// vez de emoji: emoji renderiza diferente em cada SO/navegador (as vezes
// vira 🟢 gigante, as vezes um retângulo com código, dependendo da fonte do
// sistema); um ícone SVG de uma lib fica idêntico em qualquer lugar e
// aceita `size`/`color` como qualquer outro componente React.

import { Lightbulb } from 'lucide-react';

const ITEMS = [
  { status: 'muted', pattern: 'Aceso fixo', meaning: 'Esperando leitura' },
  { status: 'normal', pattern: 'Pisca devagar', meaning: 'Normal' },
  { status: 'alert', pattern: 'Pisca rápido', meaning: 'Leitura inconsistente' },
  { status: 'error', pattern: '2 piscadas + pausa', meaning: 'Alerta — acima do limiar' },
];

export default function LedLegend() {
  return (
    <div className="card">
      <div className="card-title">
        <Lightbulb size={18} strokeWidth={2} />
        <span>O que o LED do sensor significa</span>
      </div>

      {ITEMS.map((item) => (
        <div className="legend-row" key={item.pattern}>
          <span className={`dot dot-${item.status}`} />
          <span>
            <b>{item.pattern} — </b>
            {item.meaning}
          </span>
        </div>
      ))}
    </div>
  );
}
