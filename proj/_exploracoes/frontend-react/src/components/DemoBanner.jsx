// Aviso mostrado quando a tela está usando mockData.js em vez de dados
// reais (ver useLogs.js / useConfig.js). Existe pra deixar explícito que o
// que está na tela é fake — bom pra apresentar/testar sem o backend, mas
// sem se confundir depois achando que era dado real.
import { WifiOff } from 'lucide-react';

export default function DemoBanner() {
  return (
    <div className="demo-banner">
      <WifiOff size={14} />
      <span>Modo demonstração — sem conexão com o backend</span>
    </div>
  );
}
