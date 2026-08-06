# Checklist de arrumação — pré-GitHub / LinkedIn

Lista acionável derivada da [auditoria](./CODE_REVIEW.md). Marque conforme for
fazendo. Ordenada por impacto: primeiro o que **destrava a demo pública**,
depois o que deixa o **repositório apresentável**, por fim **polimento**.

> Legenda de esforço: 🟢 rápido (min) · 🟡 médio (~1h) · 🔴 maior (algumas horas)

## 🚦 Bloco 1 — Antes de expor a URL pública (segurança/infra)

- [ ] 🟢 **Chave de API forte no Fly** — `fly secrets set API_KEY=<forte>`; usar
      a MESMA no build do Flutter (`--dart-define=API_KEY=...`). *(S1)*
- [ ] 🟢 **Travar o CORS** — `fly secrets set ALLOWED_ORIGIN=https://<seu-app>`
      e `fly deploy` depois de saber a URL do frontend. *(S3)*
- [ ] 🟡 **Decidir sobre o `DELETE /logging`** — na demo pública, remover o
      endpoint OU protegê-lo com um token de admin separado, pra ninguém zerar
      seu histórico. *(S2)*
- [ ] 🟢 **Pinar a imagem base** — no `Dockerfile`, trocar `alpine:latest` por
      `alpine:3.20` (build reprodutível). *(I1)*

## 🖼️ Bloco 2 — Deixar o repositório apresentável

- [ ] 🟢 **Adicionar `LICENSE`** (MIT) na raiz — sem isso o repo é "todos os
      direitos reservados". *(P1)*
- [ ] 🟡 **Screenshots + GIF no README** — criar `docs/img/`, colocar prints das
      telas Logs e Config e um GIF do LED do ESP32; referenciar no README. *(P2)*
- [ ] 🟡 **CI no GitHub Actions** — workflow que roda `go test`/`go build` a
      cada push, e o selo no topo do README. *(I4, P3)*
- [ ] 🟢 **Metadados no GitHub** — description + topics
      (`esp32`, `golang`, `flutter`, `iot`, `embedded`, `sqlite`) + link do demo
      no campo "Website" do repo. *(P4)*
- [ ] 🟢 **Link do demo ao vivo no topo do README.**

## ✨ Bloco 3 — Polimento (opcional, mas soma)

- [ ] 🟢 **`gofmt -w proj/backend`** — deixar os 4 arquivos Go no padrão. *(C3)*
- [ ] 🟡 **Validar payload do `POST /logging`** — rejeitar `bpm`/`spo2` fora de
      faixa plausível. *(C2)*
- [ ] 🟡 **Rate limiting** por IP nas rotas de escrita. *(S4)*
- [ ] 🟢 **`class NOT NULL DEFAULT 0`** (ou `sql.NullInt64`) no read de logs. *(C1)*
- [ ] 🟢 **Cold start** — se quiser o link sempre rápido, `min_machines_running = 1`
      no `fly.toml` (custa um pouco). *(I2)*
- [ ] 🟢 **Backup** — confirmar snapshots do volume no Fly. *(I3)*

## 📢 Bloco 4 — Publicar

- [ ] 🟢 `git push -u origin dev` e abrir PR pra `main` (`gh pr create --fill`).
- [ ] 🟢 Conferir que a demo está no ar (`curl .../healthz`) e o app abre.
- [ ] 🟢 **Post no LinkedIn:** 1 frase do problema + diagrama de arquitetura do
      README + link do demo + link do repo + stacks (ESP32 / Go / Flutter).
      Ângulo forte: *"sistema embarcado ponta a ponta, do sensor ao app na nuvem"*.

---

### Já concluído (não precisa refazer)
- ✅ Histórico git limpo, commits temáticos (era 1 commit só).
- ✅ Credenciais de WiFi/API fora do versionamento (`secrets.h`, `.env`).
- ✅ Backend endurecido: `/healthz`, graceful shutdown, WAL, limite de body,
      usuário não-root no Docker.
- ✅ Testes Go do fluxo de logging/config/auth.
- ✅ README com arquitetura, API e guia de deploy.
- ✅ Config de deploy pronta (`fly.toml`, script de build web).
- ✅ Frontends alternativos arquivados em `_exploracoes/`.
