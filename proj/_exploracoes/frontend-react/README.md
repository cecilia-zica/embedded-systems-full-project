# Monitor Cardíaco — versão React (tema claro)

Réplica do app Flutter (`proj/app`) e do protótipo `proj/web`, mas como um
projeto React "de verdade" (Vite + JSX + hooks + componentes em arquivos
separados), do jeito que você vai encontrar em tutoriais e vagas de
emprego.

Identidade visual: **navy + dourado**, tema **claro** — o par desta versão
é `proj/frontend-scala`, que usa a mesma paleta só que no tema escuro. É de
propósito: mesma lógica, dois frameworks, visual claramente diferenciável
numa demonstração lado a lado.

## Rodando

```bash
npm install
cp .env.example .env   # ajuste o IP do backend se precisar
npm run dev
```

Abre em `http://localhost:5173`. Se o backend Go (`proj/backend`) não
estiver acessível, o app cai sozinho em **modo demonstração** (dados fake,
com um aviso na tela) depois de ~4s — dá pra ver o app inteiro funcionando
sem precisar ligar nada. Ver "Modo demonstração" mais abaixo.

## Onde ler primeiro

Leia os arquivos nessa ordem — ela segue o fluxo de dados do app, do mais
simples pro mais complexo:

1. `src/config.js` — de onde vem a config (URL do backend, API key).
2. `src/theme/tokens.css` — a paleta inteira do app num lugar só.
3. `src/utils/format.js` — funções puras, sem hooks, bom pra pegar a sintaxe JS.
4. `src/api/apiService.js` e `src/api/mockData.js` — chamadas HTTP + o fallback de demonstração.
5. `src/hooks/useLogs.js` e `src/hooks/useConfig.js` — onde a lógica de estado/dados mora.
6. `src/components/` — peças pequenas e reaproveitáveis (LedLegend, StatusBadge, ConfirmDialog, DemoBanner).
7. `src/screens/LogsScreen.jsx` e `src/screens/ConfigScreen.jsx` — telas, agora só com JSX (a lógica já saiu pros hooks).
8. `src/App.jsx` e `src/main.jsx` — como tudo se conecta.

## Arquitetura — por que separado assim

- **`config.js`**: única fonte de configuração (lê `.env` via `import.meta.env`). Nenhum outro arquivo lê variável de ambiente diretamente.
- **`hooks/`**: telas (`screens/`) só desenham UI; hooks cuidam de buscar/guardar dado. Assim dá pra testar a lógica sem precisar renderizar nada, e a tela fica curta o bastante pra ler de uma vez.
- **`theme/tokens.css`**: cores/espaçamento/fonte em variáveis CSS, não hex cru espalhado pelos componentes. Trocar a cor de um "alerta" no app inteiro é editar 1 linha.
- **`utils/format.js`**: `classToStatus` devolve uma palavra (`'normal'|'alert'|'error'`), não uma cor — quem decide a cor de verdade é o CSS. Separa "o que o dado significa" de "como ele aparece na tela".

## Ícones

Emojis viraram ícones da [Lucide](https://lucide.dev) (`lucide-react`) —
mesma lib usada por boa parte do ecossistema React hoje (shadcn/ui,
Radix, etc.). Ícones SVG ficam idênticos em qualquer SO/navegador; emoji
não (renderiza diferente em cada fonte do sistema).

## Modo demonstração

Quando o fetch pro backend falha (ou dá timeout depois de 4s — ver
`REQUEST_TIMEOUT_MS` em `apiService.js`), os hooks caem pra dados fake
definidos em `src/api/mockData.js`, e a tela mostra um aviso ("Modo
demonstração — sem conexão com o backend"). Salvar/limpar continuam
funcionando nesse modo (só não persistem em lugar nenhum) — dá pra
demonstrar a interação completa sem precisar do ESP32 nem do backend rodando.

## Diferenças pra `proj/web`

`proj/web` é uma versão React sem build (React puro via `<script>`, sem
JSX) — feita de propósito pra funcionar offline numa apresentação. Esta
pasta é o oposto: JSX, hooks, componentes em arquivos separados, `npm
install` antes de rodar. Mesma lógica, mesmos endpoints — só a forma de
escrever muda.

## Diferença pro app Flutter

Esse React só tem a tela de Logs e Config (sem RFID/BLE, sem TinyML — isso
mora no firmware/backend). O objetivo é comparar os *padrões de UI*
(estado, listas, formulários, navegação por abas), não reimplementar o app
inteiro.
