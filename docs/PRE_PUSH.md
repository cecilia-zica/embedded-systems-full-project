# Checklist pré-push — revisão antes do push final

Rode as conferências, marque cada item e execute os **comandos finais** no fim.
Estado verificado nesta preparação (branch `dev`):

## 1. Higiene do repositório

- [x] Estou na branch **`dev`** — `git branch --show-current`
- [x] Working tree **limpo** — `git status --short` (nada pendente)
- [x] **Nenhum segredo real** versionado (senha/SSID do WiFi) —
      `git grep -nI -e "12345678" -e "Cecilia’s"` → vazio
- [x] `secrets.h` e `.env` **não** rastreados —
      `git ls-files | grep -E "secrets\.h$|/\.env$"` → vazio
- [x] Remote certo — `git remote get-url origin`
      → `https://github.com/cecilia-zica/embedded-systems-full-project.git`

## 2. Qualidade do código

- [x] `gofmt` limpo — `cd proj/backend && test -z "$(gofmt -l .)"`
- [x] `go vet` sem avisos — `go vet ./...`
- [x] Testes verdes — `go test ./...`
- [ ] (opcional) App Flutter compila — `cd proj/app && flutter build web` já foi
      validado antes; refaça se mexeu no app
- [ ] (opcional) Imagem Docker atual — `docker compose up --build -d` e
      `curl localhost:8080/healthz` → `{"status":"ok"}`

## 3. Apresentação

- [x] README com selos de **CI** e **licença** e seção de Demonstração
- [ ] Screenshots/GIF adicionados em `docs/img/` e bloco de imagens
      **descomentado** no README *(pendente — você tira os prints)*
- [x] `LICENSE` (MIT), `SECURITY.md` e `.env.example` presentes

## 4. Lembretes de deploy (não bloqueiam o push, são do Fly/GitHub Pages)

- [ ] `fly secrets set API_KEY=<chave-forte>` (nunca o `zica123` em produção)
- [ ] `fly secrets set ALLOWED_ORIGIN=https://<domínio-do-app>` (trava o CORS)
- [ ] Decidir o `DELETE /logging` na demo pública (ver `SECURITY.md`)
- [ ] Colar o **link do demo** no topo do README e no campo "Website" do repo

---

## Comandos finais

```bash
# 0. Confirmação rápida antes de subir
git switch dev
git status --short           # tem que estar vazio
cd proj/backend && go test ./... && cd ../..

# 1. Push da branch dev
git push -u origin dev

# 2. Abrir o Pull Request pra main (preenche com os commits)
gh pr create --base main --head dev --fill

# 3. Depois de revisar no GitHub, mergear PRESERVANDO os commits
#    (NÃO use squash — o histórico temático é o valor aqui)
gh pr merge --merge --delete-branch

# --- Alternativa sem PR (merge direto, histórico linear) ---
# git switch main && git merge --ff-only dev && git push origin main && git switch dev
```

> Depois do merge, o **CI** roda sozinho e o selo do README fica verde.
> Em seguida siga o guia de deploy do README (Fly.io + Flutter Web).
