# Política de Segurança

Este é um projeto acadêmico/portfólio. Ainda assim, seguem o modelo de ameaça
e as práticas adotadas, de forma transparente.

## Como reportar uma vulnerabilidade

Abra uma *issue* privada (Security Advisory) no GitHub ou entre em contato pelo
e-mail do perfil. Por favor, não abra issue pública para falhas sensíveis.

## Gestão de segredos

- Credenciais de WiFi e chave de API **não** ficam versionadas. Ficam em
  arquivos locais ignorados pelo git: `proj/firmware/src/secrets.h` (a partir de
  `secrets.example.h`) e `.env` do frontend. O backend lê tudo de variáveis de
  ambiente (`API_KEY`, `ALLOWED_ORIGIN`, `DB_PATH`).
- Em produção, defina uma `API_KEY` forte (ex.: `fly secrets set API_KEY=...`).
  Nunca use o default de desenvolvimento (`zica123`) exposto publicamente.

## Controles implementados

- Autenticação por `X-API-Key` comparada em **tempo constante** (`crypto/subtle`).
- **Rate limiting** por IP nas rotas de escrita (POST/DELETE).
- Limite de tamanho de corpo (1 MB) nas rotas POST.
- Respostas de erro **genéricas** (sem vazar detalhes internos do banco).
- CORS restringível por `ALLOWED_ORIGIN`.
- Container roda como usuário **não-root**.
- Timeouts de servidor e graceful shutdown.

## Limitações conhecidas (por ser demo)

- **Chave de API única e pública nos clientes.** Qualquer app web/mobile precisa
  embutir a chave, então ela chega ao usuário final. Para um produto real, o
  correto seria autenticação por dispositivo (token só do ESP32) para escrita e
  sessão/OAuth para o app, com a API pública apenas para leitura.
- **`DELETE /api/v1/logging` é destrutivo** e hoje é protegido apenas por essa
  chave. Numa exposição pública real, deveria exigir um token de administrador
  separado ou ser removido da API pública.
- Sem TLS na aplicação (delegado ao proxy da plataforma de deploy, ex.: Fly.io).

Essas limitações são aceitáveis para uma demonstração controlada, mas estão
documentadas para deixar claro o que mudaria num cenário de produção.
