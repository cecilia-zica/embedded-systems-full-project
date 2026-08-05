// Tela 2 — mostra e edita o limiar de alerta (bpm_threshold) que o ESP32
// consulta a cada 30s. Depois de extrair useConfig.js, essa função só
// desenha a UI e lê/chama o que o hook expõe.

import { useConfig } from '../hooks/useConfig';
import DemoBanner from '../components/DemoBanner';

export default function ConfigScreen() {
  const { threshold, setThreshold, alertEnabled, setAlertEnabled, loading, saving, snackbar, usingMock, save } =
    useConfig();

  if (loading) {
    return <div className="spinner" />;
  }

  return (
    <div>
      <div className="appbar">Configurações</div>
      {usingMock && <DemoBanner />}

      <div className="card">
        <label className="field-label">Limiar de alerta (BPM)</label>
        <input
          className="text-input"
          type="number"
          value={threshold}
          onChange={(ev) => setThreshold(ev.target.value)}
          placeholder="ex: 120"
        />

        <div className="switch-row">
          <span>Alerta ativado</span>
          {/* Não existe <Switch> pronto em HTML puro — esse botão é um
              "switch" desenhado via CSS (ver .switch em styles.css). */}
          <button
            className={`switch ${alertEnabled ? 'on' : ''}`}
            onClick={() => setAlertEnabled(!alertEnabled)}
            aria-pressed={alertEnabled}
          >
            <span className="switch-thumb" />
          </button>
        </div>
      </div>

      <button className="btn-primary" disabled={saving} onClick={save}>
        {saving ? 'Salvando...' : 'SALVAR'}
      </button>

      {snackbar && <div className="snackbar">{snackbar}</div>}
    </div>
  );
}
