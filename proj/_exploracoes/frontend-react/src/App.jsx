// Componente raiz — equivalente a MyApp + HomeNav + _HomeNavState no Dart:
// guarda qual aba está ativa e desenha a bottom nav.

import { useState } from 'react';
import { ClipboardList, Settings } from 'lucide-react';
import LogsScreen from './screens/LogsScreen';
import ConfigScreen from './screens/ConfigScreen';

export default function App() {
  const [tab, setTab] = useState(0);

  return (
    <div style={{ display: 'flex', flexDirection: 'column', minHeight: '100vh' }}>
      <div style={{ flex: 1 }}>{tab === 0 ? <LogsScreen /> : <ConfigScreen />}</div>

      <div className="bottom-nav">
        <button className={`nav-item ${tab === 0 ? 'active' : ''}`} onClick={() => setTab(0)}>
          <ClipboardList size={20} />
          <span>Logs</span>
        </button>
        <button className={`nav-item ${tab === 1 ? 'active' : ''}`} onClick={() => setTab(1)}>
          <Settings size={20} />
          <span>Config</span>
        </button>
      </div>
    </div>
  );
}
