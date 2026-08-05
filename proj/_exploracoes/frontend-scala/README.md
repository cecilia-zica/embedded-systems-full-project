# Monitor Cardíaco — versão Scala.js + Laminar (tema escuro)

Réplica do app Flutter (`proj/app`) e das versões React (`proj/web`,
`proj/frontend-react`), agora em Scala compilado pra JavaScript via
Scala.js, usando a lib de UI reativa Laminar.

Identidade visual: **navy + dourado**, tema **escuro** — o par desta
versão é `proj/frontend-react`, que usa a mesma paleta só que no tema
claro. É de propósito: mesma lógica, dois frameworks, visual claramente
diferenciável numa demonstração lado a lado (e navy+dourado combina
particularmente bem em fundo escuro).

Scala não é uma escolha comum pra frontend web — o normal é rodar Scala só
no backend (JVM). Essa versão existe como exercício comparativo: mesma
lógica, três jeitos bem diferentes de organizar o mesmo problema
(Flutter/widgets, React/hooks, Laminar/reativo).

## Pré-requisitos

Diferente do React (só precisa de Node), Scala.js precisa de:

- **JDK** (Java 11+) — o sbt roda sobre a JVM, mesmo compilando pra JS no final.
- **sbt** — a ferramenta de build do Scala (equivalente ao `npm`/`pubspec.yaml`).

No Mac, o jeito mais simples: `brew install sbt` (já baixa um JDK junto).
Se o `brew install sbt` ficar preso compilando um monte de dependência
gráfica (cairo/glib/gtk) do zero, é porque não existe binário pronto de
`openjdk` pra sua versão do macOS — nesse caso baixe um JDK portátil (ex:
[Temurin](https://adoptium.net), `.tar.gz`, sem precisar de instalador/sudo)
e rode `brew install sbt --ignore-dependencies` com esse JDK no `JAVA_HOME`.

## Rodando

```bash
# compila o Scala pra um arquivo JavaScript (fica em target/.../main.js)
sbt fastLinkJS

# sobe um servidor estático qualquer nessa pasta pra abrir o index.html
npx serve .
# (ou: python3 -m http.server 8000)
```

Depois abre a URL que o servidor mostrar. Se o backend Go (`proj/backend`)
não estiver acessível, o app cai sozinho em **modo demonstração** (dados
fake, com um aviso na tela) depois de ~4s — ver "Modo demonstração" mais
abaixo. Config do backend fica em `index.html` (`window.__CONFIG__`), lida
por `Env.scala`.

Pra ficar recompilando sozinho a cada mudança (tipo o hot-reload do Vite):

```bash
sbt ~fastLinkJS
```

## Estrutura dos pacotes

```
monitor/
├── api/          Config (Env), chamadas HTTP (Api), dados de demonstração (MockData)
├── model/        case classes que espelham o JSON do backend (LogEntry, Config)
├── format/       funções puras de formatação (Format)
├── icons/        ícones SVG da Lucide embutidos como código Scala (Icons)
├── components/   peças pequenas e reaproveitáveis (LedLegend, StatusBadge, ConfirmDialog, DemoBanner)
├── screens/      as duas telas (LogsScreen, ConfigScreen)
├── App.scala     raiz — navegação por abas
└── Main.scala    entry point
```

Mesma ideia de separação por responsabilidade que a versão React
(`config.js` → `api/Env.scala`, `hooks/` → a lógica dentro de cada
`screens/*.scala`, `theme/tokens.css` → variáveis CSS em `styles.css`).

## Onde ler primeiro

1. `src/main/scala/monitor/model/Models.scala` — como o JSON do backend
   vira tipos Scala (compare com o Map<String,dynamic> solto do Dart/JS).
2. `src/main/scala/monitor/format/Format.scala` — funções puras, bom pra
   pegar a sintaxe Scala sem se preocupar com UI ainda.
3. `src/main/scala/monitor/api/Api.scala` e `MockData.scala` — chamadas
   HTTP com Future + o fallback de demonstração.
4. `src/main/scala/monitor/icons/Icons.scala` — ícones SVG (Lucide)
   embutidos como código, sem depender de npm.
5. `src/main/scala/monitor/components/` — LedLegend (sem estado),
   ConfirmDialog/StatusBadge/DemoBanner (recebem dado por parâmetro).
6. `src/main/scala/monitor/screens/LogsScreen.scala` e `ConfigScreen.scala`
   — onde `Var`/`Signal` entram (o "estado reativo" do Laminar).
7. `src/main/scala/monitor/App.scala` e `Main.scala` — como tudo se conecta.

## Ícones sem npm

Os mesmos ícones da versão React (lib [Lucide](https://lucide.dev)), só
que aqui copiados diretamente como coordenadas SVG em `Icons.scala` — esse
projeto não usa nenhum bundler (só sbt puro), então não dá pra fazer
`npm install lucide`. Resultado visual idêntico, forma de "instalar"
diferente.

## Modo demonstração

Mesma ideia da versão React: quando o fetch pro backend falha (ou dá
timeout depois de 4s — ver `requestTimeoutMs` em `Api.scala`), as telas
caem pros dados fake de `MockData.scala` e mostram um aviso ("Modo
demonstração — sem conexão com o backend"). Salvar/limpar continuam
funcionando nesse modo (só não persistem em lugar nenhum).

## React vs Laminar — a diferença que mais importa

No React (`useState`), mudar o estado faz o componente inteiro **rodar de
novo** (a função inteira reexecuta) e o Virtual DOM calcula o que mudou.

Em Laminar, `Var`/`Signal` são um fluxo de dados: você "amarra" um pedaço
específico do DOM a um Signal com `<--`, e só AQUELE pedaço é atualizado
quando o valor muda — o resto do componente nunca reexecuta. É mais parecido
com o `setState(() { _campo = valor })` do Flutter granular, só que sem
precisar redesenhar o `build()` inteiro.

Nenhum dos dois é "melhor" — são modelos mentais diferentes pro mesmo
problema (estado muda → tela atualiza). Vale comparar `LogsScreen.jsx` com
`LogsScreen.scala` lado a lado pra sentir a diferença.
