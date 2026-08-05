// Tela 1 — mostra o histórico (logs) de leituras que o ESP32 mandou.
// Equivalente a LogsScreen + _LogsScreenState no Dart. Depois de extrair
// useLogs.js, essa função ficou só com a parte de "desenhar a UI" — a
// lógica de buscar/apagar dados mora inteira no hook.

import { useState } from 'react';
import { RefreshCw, Trash2, HeartPulse } from 'lucide-react';
import { useLogs } from '../hooks/useLogs';
import { classToStatus, formatTimestamp } from '../utils/format';
import LedLegend from '../components/LedLegend';
import ConfirmDialog from '../components/ConfirmDialog';
import StatusBadge from '../components/StatusBadge';
import DemoBanner from '../components/DemoBanner';

export default function LogsScreen() {
  const { logs, loading, error, usingMock, refresh, remove } = useLogs();
  const [confirmingClear, setConfirmingClear] = useState(false);

  async function confirmClear() {
    setConfirmingClear(false);
    await remove();
  }

  let body;
  if (loading) {
    body = <div className="spinner" />;
  } else if (error) {
    body = <div className="center-msg">Erro: {error}</div>;
  } else if (logs.length === 0) {
    body = <div className="center-msg">Nenhuma leitura ainda</div>;
  } else {
    body = (
      <div>
        {logs.map((log) => {
          const status = classToStatus(log.class);
          return (
            <div className="log-item" key={log.id ?? `${log.user_id}-${log.created_at}`}>
              <div className={`avatar status-${status}`}>
                <HeartPulse size={20} strokeWidth={2} />
              </div>
              <div className="log-main">
                <div className="log-title">
                  {log.bpm} BPM · SpO2 {log.spo2}%
                </div>
                <div className="log-subtitle">
                  {log.user_id} — {formatTimestamp(log.created_at)}
                </div>
              </div>
              <StatusBadge status={status} />
            </div>
          );
        })}
      </div>
    );
  }

  return (
    <div style={{ display: 'flex', flexDirection: 'column', height: '100%' }}>
      <div className="appbar">
        Logs
        <div className="appbar-actions">
          <button className="icon-btn" onClick={refresh} title="Atualizar">
            <RefreshCw size={18} />
          </button>
          <button className="icon-btn" onClick={() => setConfirmingClear(true)} title="Limpar histórico">
            <Trash2 size={18} />
          </button>
        </div>
      </div>

      <div className="content">
        {usingMock && <DemoBanner />}
        <LedLegend />
        {body}
      </div>

      <ConfirmDialog
        open={confirmingClear}
        title="Limpar log?"
        message="Apaga todo o histórico de leituras. Não dá pra desfazer."
        confirmLabel="Limpar"
        onCancel={() => setConfirmingClear(false)}
        onConfirm={confirmClear}
      />
    </div>
  );
}
