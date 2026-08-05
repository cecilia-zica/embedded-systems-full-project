# Explorações — mesmo app, outras stacks

O frontend oficial do projeto (e o alvo de deploy) é o app **Flutter** em
[`../app`](../app). Estas são implementações alternativas do **mesmo** frontend,
consumindo a **mesma** API Go, feitas para explorar diferentes stacks. Não
entram no deploy, ficam aqui como referência.

| Pasta | Stack | Observação |
|-------|-------|------------|
| `frontend-react/` | React + Vite | Versão mais completa; hooks, tema, componentes. |
| `frontend-scala/` | Scala.js | Porte da mesma UI para Scala.js. |
| `web/` | HTML + JS puro (React via CDN) | Protótipo sem build step. |

Cada uma tem o próprio README/instruções de execução. Todas apontam para o
backend em [`../backend`](../backend) e usam a mesma API key de dev.
