// Versão React (web, sem build/npm) do mesmo app Flutter — só pra testar no
// celular na apresentação. Mesmos endpoints, mesma lógica, mesmos padrões.
// Sem JSX/Babel de propósito: só React puro via CDN local, funciona abrindo
// o HTML direto, sem precisar compilar nada.

const e = React.createElement; // atalho pra não escrever React.createElement toda hora

// ---------- ApiService — equivalente ao services/api_service.dart ----------
const ApiService = {
  // mesmo IP/porta do backend que o app Flutter já usa
  baseUrl: 'http://172.20.10.2:8080',
  apiKey: 'zica123', // tem que bater com defaultAPIKey em proj/backend/middleware.go

  headers: {
    'Content-Type': 'application/json',
    'X-API-Key': 'zica123',
  },

  // GET /api/v1/logging — busca o histórico de leituras
  async getLogs() {
    const res = await fetch(`${this.baseUrl}/api/v1/logging`, { headers: this.headers });
    if (!res.ok) throw new Error(`erro ao buscar logs: ${res.status}`);
    return res.json(); // fetch().json() == json.decode() do Dart
  },

  // GET /api/v1/controle — busca a config atual
  async getConfig() {
    const res = await fetch(`${this.baseUrl}/api/v1/controle`, { headers: this.headers });
    if (!res.ok) throw new Error(`erro ao buscar config: ${res.status}`);
    return res.json();
  },

  // POST /api/v1/controle — salva threshold/alerta
  async postConfig(threshold, alertEnabled) {
    const res = await fetch(`${this.baseUrl}/api/v1/controle`, {
      method: 'POST',
      headers: this.headers,
      body: JSON.stringify({ bpm_threshold: threshold, alert_enabled: alertEnabled }),
    });
    if (!res.ok) throw new Error(`erro ao salvar config: ${res.status}`);
  },
};

// ---------- helpers (mesma lógica de logs_screen.dart) ----------
function classToText(c) {
  if (c === 0) return 'Normal';
  if (c === 1) return 'Alerta';
  return 'Erro Leitura';
}

function classToColor(c) {
  if (c === 0) return '#2e7d32'; // verde
  if (c === 1) return '#ef6c00'; // laranja
  return '#c62828'; // vermelho
}

// new Date(isoString) já entende o "Z" (UTC) do backend e getHours()/etc já
// devolvem hora local do navegador — o JS faz o mesmo que o .toLocal() do Dart
// automaticamente, sem precisar converter nada na mão
function formatTimestamp(iso) {
  const dt = new Date(iso);
  const two = (n) => String(n).padStart(2, '0');
  return `${two(dt.getDate())}/${two(dt.getMonth() + 1)} ${two(dt.getHours())}:${two(dt.getMinutes())}:${two(dt.getSeconds())}`;
}

// ---------- LedLegend — equivalente a _buildLedLegend() ----------
function LedLegend() {
  const items = [
    { color: '#9e9e9e', pattern: 'Aceso fixo', meaning: 'Esperando leitura' },
    { color: '#2e7d32', pattern: 'Pisca devagar', meaning: 'Normal' },
    { color: '#ef6c00', pattern: 'Pisca rápido', meaning: 'Leitura inconsistente' },
    { color: '#c62828', pattern: '2 piscadas + pausa', meaning: 'Alerta — acima do limiar' },
  ];

  return e('div', { className: 'card' }, [
    e('div', { key: 'title', style: { display: 'flex', alignItems: 'center', gap: 8, marginBottom: 8, fontWeight: 600 } }, [
      e('span', { key: 'icon' }, '💡'),
      e('span', { key: 'text' }, 'O que o LED do sensor significa'),
    ]),
    ...items.map((item) =>
      e('div', { className: 'legend-row', key: item.pattern }, [
        e('span', { className: 'dot', style: { background: item.color }, key: 'dot' }),
        e('span', { key: 'text' }, [
          e('b', { key: 'b' }, `${item.pattern} — `),
          item.meaning,
        ]),
      ])
    ),
  ]);
}

// ---------- LogsScreen — equivalente a LogsScreen + _LogsScreenState ----------
function LogsScreen() {
  // useState = a mesma ideia das variáveis dentro de _LogsScreenState (_logs,
  // _loading, _erro); cada useState é uma variável de estado independente
  const [logs, setLogs] = React.useState([]);
  const [loading, setLoading] = React.useState(true);
  const [erro, setErro] = React.useState(null);

  async function fetchLogs() {
    try {
      const data = await ApiService.getLogs();
      // no React não existe setState({...}) genérico como no Flutter — cada
      // useState tem sua própria função de atualização, mas o efeito é o mesmo:
      // muda o dado E agenda o redesenho, os dois numa chamada só
      setLogs(data);
      setLoading(false);
      setErro(null);
    } catch (err) {
      setErro(String(err));
      setLoading(false);
    }
  }

  // useEffect(() => {...}, []) com array vazio = roda 1 vez só, quando o
  // componente "nasce" — é o equivalente direto do initState() do Flutter
  React.useEffect(() => {
    fetchLogs();
  }, []);

  let body;
  if (loading) {
    body = e('div', { className: 'spinner' });
  } else if (erro) {
    body = e('div', { className: 'center-msg' }, `Erro: ${erro}`);
  } else if (logs.length === 0) {
    body = e('div', { className: 'center-msg' }, 'Nenhuma leitura ainda');
  } else {
    body = e('div', null, logs.map((log) =>
      e('div', { className: 'log-item', key: log.id ?? Math.random() }, [
        e('div', {
          className: 'avatar',
          key: 'avatar',
          style: { background: classToColor(log.class) + '26' }, // 26 = ~15% de opacidade em hex
        }, '❤️'),
        e('div', { className: 'log-main', key: 'main' }, [
          e('div', { className: 'log-title', key: 'title' }, `${log.bpm} BPM   ·   SpO2 ${log.spo2}%`),
          e('div', { className: 'log-subtitle', key: 'subtitle' }, `${log.user_id} — ${formatTimestamp(log.created_at)}`),
        ]),
        e('span', {
          className: 'chip',
          key: 'chip',
          style: { background: classToColor(log.class) + '26', color: classToColor(log.class) },
        }, classToText(log.class)),
      ])
    ));
  }

  return e('div', { style: { display: 'flex', flexDirection: 'column', height: '100%' } }, [
    e('div', { className: 'appbar', key: 'appbar' }, [
      'Logs ',
      // botão de atualizar manual — equivalente ao IconButton(refresh) do AppBar
      e('button', {
        key: 'refresh',
        onClick: fetchLogs,
        style: { border: 'none', background: 'none', fontSize: 16, cursor: 'pointer', float: 'right' },
      }, '🔄'),
    ]),
    e('div', { className: 'content', key: 'content' }, [
      e(LedLegend, { key: 'legend' }),
      body,
    ]),
  ]);
}

// ---------- ConfigScreen — equivalente a ConfigScreen + _ConfigScreenState ----------
function ConfigScreen() {
  const [threshold, setThreshold] = React.useState('');
  const [alertEnabled, setAlertEnabled] = React.useState(true);
  const [loading, setLoading] = React.useState(true);
  const [saving, setSaving] = React.useState(false);
  const [snackbar, setSnackbar] = React.useState(null); // mensagem temporária, tipo SnackBar

  async function loadCurrentConfig() {
    try {
      const config = await ApiService.getConfig();
      setThreshold(String(config.bpm_threshold));
      setAlertEnabled(config.alert_enabled);
      setLoading(false);
    } catch (err) {
      setLoading(false);
      showSnackbar(`Erro ao carregar configuração: ${err}`);
    }
  }

  React.useEffect(() => {
    loadCurrentConfig();
  }, []);

  function showSnackbar(msg) {
    setSnackbar(msg);
    setTimeout(() => setSnackbar(null), 3000); // some sozinho, igual o SnackBar do Flutter
  }

  async function saveConfig() {
    // Number.isInteger + parseInt não lança exceção como int.parse — mesma
    // ideia do int.tryParse do Dart: valida antes de gastar uma chamada de rede
    const parsed = parseInt(threshold, 10);
    if (Number.isNaN(parsed)) {
      showSnackbar('Digite um número válido pro limiar de BPM');
      return;
    }

    setSaving(true);
    try {
      await ApiService.postConfig(parsed, alertEnabled);
      showSnackbar('Configuração salva!');
    } catch (err) {
      showSnackbar(`Erro ao salvar: ${err}`);
    } finally {
      setSaving(false);
    }
  }

  if (loading) {
    return e('div', { className: 'spinner' });
  }

  return e('div', null, [
    e('div', { className: 'appbar', key: 'appbar' }, 'Configurações'),
    e('div', { className: 'card', key: 'card' }, [
      e('label', { className: 'field-label', key: 'label' }, 'Limiar de alerta (BPM)'),
      e('input', {
        key: 'input',
        className: 'text-input',
        type: 'number',
        value: threshold,
        // onChange = onChanged do TextField no Flutter: dispara a cada tecla
        onChange: (ev) => setThreshold(ev.target.value),
        placeholder: 'ex: 120',
      }),
      e('div', { className: 'switch-row', key: 'switch-row' }, [
        e('span', { key: 'label' }, 'Alerta ativado'),
        e('button', {
          key: 'switch',
          className: 'switch',
          style: {
            background: alertEnabled ? '#c62828' : '#bbb',
          },
          onClick: () => setAlertEnabled(!alertEnabled),
        }, e('span', {
          style: { position: 'absolute', left: alertEnabled ? 23 : 3, top: 3, width: 20, height: 20, borderRadius: '50%', background: '#fff', transition: 'left 0.15s' },
        })),
      ]),
    ]),
    e('button', {
      key: 'save',
      className: 'btn-primary',
      disabled: saving,
      onClick: saveConfig,
    }, saving ? 'Salvando...' : 'SALVAR'),
    snackbar && e('div', { className: 'snackbar', key: 'snackbar' }, snackbar),
  ]);
}

// ---------- App — equivalente a MyApp + HomeNav + _HomeNavState ----------
function App() {
  // useState(0) = igual o int _indiceAtual = 0 do Flutter (a gente já viu
  // esse mesmo comentário no HomeNav) — guarda qual aba está selecionada
  const [tab, setTab] = React.useState(0);

  return e('div', { style: { display: 'flex', flexDirection: 'column', minHeight: '100vh' } }, [
    e('div', { style: { flex: 1 }, key: 'body' }, tab === 0 ? e(LogsScreen) : e(ConfigScreen)),
    e('div', { className: 'bottom-nav', key: 'nav' }, [
      e('button', {
        key: 'logs',
        className: `nav-item ${tab === 0 ? 'active' : ''}`,
        // onClick(0) = onTap: (index) => setState(() => _indiceAtual = index) do Flutter
        onClick: () => setTab(0),
      }, [e('span', { className: 'icon', key: 'i' }, '📋'), e('span', { key: 't' }, 'Logs')]),
      e('button', {
        key: 'config',
        className: `nav-item ${tab === 1 ? 'active' : ''}`,
        onClick: () => setTab(1),
      }, [e('span', { className: 'icon', key: 'i' }, '⚙️'), e('span', { key: 't' }, 'Config')]),
    ]),
  ]);
}

// ReactDOM.createRoot(...).render(...) = o runApp(const MyApp()) do Flutter,
// "liga" a árvore de componentes na tela
const root = ReactDOM.createRoot(document.getElementById('root'));
root.render(e(App));
