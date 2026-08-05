// Custom hook — mesma ideia do useLogs.js, mas pra tela de Configurações:
// carrega o threshold/alerta atual, guarda o rascunho que o usuário está
// editando, e sabe salvar. ConfigScreen.jsx fica só com JSX depois disso.

import { useState, useEffect, useCallback } from 'react';
import { getConfig, postConfig } from '../api/apiService';
import { mockConfig } from '../api/mockData';

export function useConfig() {
  const [threshold, setThreshold] = useState('');
  const [alertEnabled, setAlertEnabled] = useState(true);
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [snackbar, setSnackbar] = useState(null); // mensagem temporária, tipo o SnackBar do Flutter
  const [usingMock, setUsingMock] = useState(false);

  const showSnackbar = useCallback((msg) => {
    setSnackbar(msg);
    setTimeout(() => setSnackbar(null), 3000); // some sozinho, igual o SnackBar do Flutter
  }, []);

  useEffect(() => {
    getConfig()
      .then((config) => {
        setThreshold(String(config.bpm_threshold));
        setAlertEnabled(config.alert_enabled);
        setUsingMock(false);
      })
      .catch(() => {
        // backend fora do ar — pré-preenche com a config de demonstração em
        // vez de deixar o formulário vazio (ver mockData.js).
        setThreshold(String(mockConfig.bpm_threshold));
        setAlertEnabled(mockConfig.alert_enabled);
        setUsingMock(true);
      })
      .finally(() => setLoading(false));
  }, []);

  const save = useCallback(async () => {
    // Number.parseInt não lança exceção quando falha (devolve NaN) — mesma
    // ideia do int.tryParse do Dart: valida ANTES de gastar uma chamada de rede.
    const parsed = Number.parseInt(threshold, 10);
    if (Number.isNaN(parsed)) {
      showSnackbar('Digite um número válido pro limiar de BPM');
      return;
    }

    setSaving(true);

    // sem backend de verdade pra salvar em modo demonstração — só simula um
    // delay de rede, pra dar pra testar a interação mesmo offline.
    if (usingMock) {
      setTimeout(() => {
        setSaving(false);
        showSnackbar('Configuração salva! (modo demonstração)');
      }, 400);
      return;
    }

    try {
      await postConfig(parsed, alertEnabled);
      showSnackbar('Configuração salva!');
    } catch (err) {
      showSnackbar(`Erro ao salvar: ${err}`);
    } finally {
      setSaving(false);
    }
  }, [threshold, alertEnabled, usingMock, showSnackbar]);

  return { threshold, setThreshold, alertEnabled, setAlertEnabled, loading, saving, snackbar, usingMock, save };
}
