// Componente reutilizável: transforma um status semântico ('normal' |
// 'alert' | 'error') num "chip" colorido com ícone + texto. Antes essa
// lógica (cor + texto) ficava espalhada dentro de LogsScreen.jsx; virar um
// componente próprio evita duplicar JSX se um dia mais de uma tela precisar
// mostrar o mesmo tipo de badge (regra de ouro do clean code: lógica usada
// em mais de 1 lugar vira uma coisa só, não é copiada).

import { CheckCircle2, AlertTriangle, XCircle } from 'lucide-react';

const CONFIG = {
  normal: { label: 'Normal', Icon: CheckCircle2 },
  alert: { label: 'Alerta', Icon: AlertTriangle },
  error: { label: 'Erro Leitura', Icon: XCircle },
};

export default function StatusBadge({ status }) {
  const { label, Icon } = CONFIG[status];
  return (
    <span className={`chip status-${status}`}>
      <Icon size={13} strokeWidth={2.5} />
      {label}
    </span>
  );
}
