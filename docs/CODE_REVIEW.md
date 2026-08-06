# Auditoria pré-lançamento — Código, Segurança e Infraestrutura

> Revisão de **máximo esforço** do estado atual (branch `dev`), **sem alterar
> código-fonte**. Foco: correção, segurança, deploy/infra e prontidão para
> publicar no GitHub / LinkedIn. Cada achado tem severidade, local e evidência.

## Evidência de verificação (executada nesta auditoria)

| Verificação | Resultado |
|---|---|
| `go vet ./...` | limpo |
| `go test ./...` | **6/6 passam** (`ok backend`) |
| `go build` (CGO off) | OK |
| `go mod verify` | `all modules verified` |
| `gofmt -l` | 4 arquivos fora do padrão (ver C3) |
| Container Docker | `Up 5h (healthy)` — healthcheck estável |
| `flutter build web` | bundle gerado (`main.dart.js` 2.0 MB) |
| `pio run` (firmware) | `firmware.bin` gerado |
| Segredos no git | **nenhuma** credencial de WiFi/SSID versionada |

**Resumo:** o código está sólido e todos os artefatos compilam. Os problemas
graves da primeira revisão (senha de WiFi hardcoded, backend sem hardening,
histórico git vazio) **já foram corrigidos**. O que resta abaixo são itens de
robustez, segurança de exposição pública e polimento — nada bloqueante para a
demo, mas relevante antes de expor uma URL pública.

---

## 1. Correção

- **C1 — [Baixo] `class` NULL é silenciosamente descartado na listagem.**
  `proj/backend/logging.go` — `handleGetLogging` faz `Scan(&l.Class)` para um
  `int` não-nulo. A coluna `class` aceita NULL (schema em `db.go`). Se algum dia
  entrar uma linha com `class` NULL, o `Scan` falha e a linha é pulada.
  *Alcance real:* o firmware sempre envia `class`, então hoje é inatingível.
  *Correção:* usar `sql.NullInt64` ou tornar a coluna `NOT NULL DEFAULT 0`.

- **C2 — [Baixo] `POST /api/v1/logging` não valida o payload.**
  `proj/backend/logging.go` — aceita `bpm`/`spo2` negativos ou absurdos e
  qualquer `user_id`. Como a API é pública (ver S1), aceita lixo.
  *Correção:* validar faixas plausíveis (ex.: `0 < bpm < 300`).

- **C3 — [Info] 4 arquivos Go não estão `gofmt`-limpos.**
  `controle.go`, `logging.go`, `main.go`, `middleware.go` — apenas espaço após
  `//` (`//comentário` → `// comentário`). Cosmético, mas um repo Go público é
  esperado ser `gofmt`-limpo. *Correção:* `gofmt -w proj/backend`.

---

## 2. Segurança

- **S1 — [Médio] Chave de API única, fraca por padrão e pública nos clientes.**
  `proj/backend/middleware.go:17` — `const defaultAPIKey = "zica123"`. Em
  qualquer cliente web/mobile a chave **necessariamente** chega ao usuário
  (fica embutida no `main.dart.js` do Flutter Web). Ou seja: quem abrir o site
  consegue ler a chave e chamar a API. O compare em tempo constante
  (`middleware.go:54`) está correto, mas não muda o modelo de exposição.
  *Mitigação já feita:* chave forte via `API_KEY` no Fly. *Recomendado p/ real:*
  escrita autenticada por dispositivo (token do ESP32) e leitura pública
  read-only, ou um proxy.

- **S2 — [Médio] `DELETE /api/v1/logging` apaga TUDO, só com a chave pública.**
  `proj/backend/logging.go` — `handleDeleteLogging` roda `DELETE FROM logs`
  sem confirmação no servidor. Com a chave sendo pública (S1), um visitante da
  demo poderia zerar o histórico. *Correção:* restringir DELETE a um token de
  admin separado, ou remover o endpoint em produção.

- **S3 — [Baixo/Médio] CORS default `*`.**
  `proj/backend/middleware.go:71` — `Access-Control-Allow-Origin: *` quando
  `ALLOWED_ORIGIN` não está setada. Correto para dev; em produção **precisa**
  ser o domínio do app. Não é bug (é configurável por env), mas é fácil de
  esquecer no deploy → vira item de checklist.

- **S4 — [Baixo] Sem rate limiting.**
  Nenhuma rota tem throttle. Com a chave pública, dá pra floodar
  `POST /logging` (crescimento do banco) ou repetir `DELETE`.
  *Correção:* limitador simples por IP (ex.: `golang.org/x/time/rate`).

- **S5 — [Bom] Pontos fortes confirmados.**
  Compare de chave em tempo constante; credenciais externalizadas em
  `secrets.h`/`.env` git-ignorados; **nenhum segredo real no histórico** (WiFi
  SSID/senha não aparecem em nenhum arquivo versionado); `go mod verify` OK.

---

## 3. Infraestrutura & Deploy

- **I1 — [Médio] Imagem base de runtime não pinada.**
  `proj/backend/Dockerfile:12` — `FROM alpine:latest`. `latest` torna o build
  não-reprodutível (a imagem muda com o tempo). *Correção:* pinar (ex.:
  `alpine:3.20`), idealmente por digest.

- **I2 — [Baixo/Médio] Cold start na demo do Fly.**
  `proj/backend/fly.toml` — `auto_stop_machines = "stop"` e
  `min_machines_running = 0`: a máquina dorme quando ociosa e a **primeira**
  requisição depois de um tempo tem latência de alguns segundos. Para um link
  de LinkedIn, o primeiro visitante pode pegar um carregamento lento.
  *Trade-off:* `min_machines_running = 1` deixa sempre-ligado (custa um pouco).

- **I3 — [Baixo] SQLite em volume único, sem backup explícito.**
  Um volume só = ponto único de falha. O Fly faz snapshots diários por padrão,
  mas vale confirmar/mencionar a estratégia. Para demo é aceitável.

- **I4 — [Baixo] Sem CI.**
  Os testes não rodam automaticamente a cada push. Um GitHub Actions com
  `go test` + `go build` adiciona um selo verde (ótimo para portfólio) e pega
  regressões. *Correção:* workflow em `.github/workflows/`.

- **I5 — [Info] Healthcheck duplicado.**
  Definido no `Dockerfile` e no `fly.toml`. Redundante, mas inofensivo.

---

## 4. Prontidão para GitHub / LinkedIn

- **P1 — [Médio] Falta `LICENSE`.** Repo público sem licença = "todos os
  direitos reservados" (ninguém pode reutilizar legalmente). Adicionar MIT.
- **P2 — [Médio] Sem screenshots/GIF no README.** `docs/img/` não existe. É o
  item que mais pesa para quem olha o repositório de fora.
- **P3 — [Baixo] Sem selo de CI** (depende de I4).
- **P4 — [Baixo] Metadados do repo no GitHub:** description, topics
  (`esp32`, `golang`, `flutter`, `iot`, `embedded`) e link do demo ao vivo.
- **P5 — [Bom] README já forte:** diagrama de arquitetura, referência de API,
  guia de deploy e nota honesta de segurança já presentes.

---

## Prioridade sugerida

1. **Antes de expor URL pública:** S3 (setar `ALLOWED_ORIGIN`), S1/S2 (chave
   forte no Fly + decidir sobre o DELETE), I1 (pinar imagem).
2. **Para o repositório ficar apresentável:** P1 (LICENSE), P2 (screenshots),
   I4/P3 (CI + selo).
3. **Polimento:** C3 (`gofmt`), C1/C2 (validação), S4 (rate limit), I2 (cold
   start), I3 (backup).

O checklist acionável está em [`PRE_LAUNCH_CHECKLIST.md`](./PRE_LAUNCH_CHECKLIST.md).
