// Custom hook — extrai TODA a lógica de estado/dados que antes vivia dentro
// de LogsScreen.jsx. É um padrão central do React "moderno": telas (screens)
// só cuidam de desenhar a UI; hooks cuidam de buscar/guardar dados. Isso é o
// equivalente, em espírito, ao ViewModel/Controller separado do Widget no
// Flutter — só que em vez de uma classe, é uma função que devolve o estado
// + as ações que mexem nele.
//
// Qualquer componente que chamar `useLogs()` ganha os mesmos logs, loading,
// erro e as mesmas funções refresh/remove — sem duplicar useState em cada
// lugar que precisar dessa lista.

import { useState, useEffect, useRef, useCallback } from 'react';
import { getLogs, deleteLogs } from '../api/apiService';
import { mockLogs } from '../api/mockData';

export function useLogs() {
  const [logs, setLogs] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  // true quando o backend não respondeu e a tela está mostrando mockData em
  // vez de dados reais — ver mockData.js. A UI usa isso pra exibir o aviso
  // de "modo demonstração" (DemoBanner) em vez de fingir que são dados reais.
  const [usingMock, setUsingMock] = useState(false);

  // useRef guarda um valor que sobrevive entre renderizações mas que, ao
  // mudar, NÃO redesenha nada — é só uma trava interna contra requisições
  // concorrentes (refresh manual clicado 2x rápido, por exemplo).
  const fetching = useRef(false);

  // useCallback memoriza a função entre renderizações: sem isso, toda vez
  // que o componente redesenha, `refresh` seria uma função NOVA — o que
  // quebraria o useEffect de quem depender dela na lista de dependências.
  const refresh = useCallback(async () => {
    if (fetching.current) return;
    fetching.current = true;
    try {
      const data = await getLogs();
      setLogs(data);
      setError(null);
      setUsingMock(false);
    } catch (err) {
      // Backend fora do ar (ou IP errado, ou você só quer ver o app sem
      // ligar nada) — em vez de deixar a tela em branco com uma mensagem de
      // erro, cai pros dados de demonstração. Não é um "erro escondido": o
      // usingMock fica true e a tela mostra um aviso deixando claro que os
      // dados são fake.
      setLogs(mockLogs);
      setError(null);
      setUsingMock(true);
    } finally {
      setLoading(false);
      fetching.current = false;
    }
  }, []);

  const remove = useCallback(async () => {
    // em modo demonstração não existe backend pra chamar — só limpa local.
    if (usingMock) {
      setLogs([]);
      return;
    }
    try {
      await deleteLogs();
      setLogs([]);
    } catch (err) {
      setError(String(err));
    }
  }, [usingMock]);

  // roda uma vez, quando quem usar esse hook "nasce" — equivalente ao
  // initState() do Flutter.
  useEffect(() => {
    refresh();
  }, [refresh]);

  return { logs, loading, error, usingMock, refresh, remove };
}
