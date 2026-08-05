// Entry point — igual ao `void main() => runApp(const MyApp());` do Flutter,
// e igual ao `root.render(e(App))` que o proj/web fazia na mão. Aqui é a
// versão "com build": import de verdade em vez de <script> globais.

import React from 'react';
import ReactDOM from 'react-dom/client';
import App from './App.jsx';
import './styles.css';

// createRoot(...).render(...) "liga" a árvore de componentes React na div
// #root do index.html — mesmo papel do runApp() do Flutter.
ReactDOM.createRoot(document.getElementById('root')).render(
  // StrictMode: só existe em desenvolvimento, ajuda a pegar bugs comuns
  // (tipo efeitos colaterais escritos errado) rodando alguns hooks 2x de
  // propósito. Não afeta o build de produção nem o usuário final.
  <React.StrictMode>
    <App />
  </React.StrictMode>,
);
