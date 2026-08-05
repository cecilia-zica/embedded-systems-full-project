// Fonte única de configuração do app — nada mais em nenhum outro arquivo
// deve ler `import.meta.env` diretamente. Se amanhã a gente trocar de Vite
// pra outra ferramenta, só esse arquivo muda.
//
// import.meta.env.VITE_* é como o Vite expõe variáveis do .env pro código
// do navegador (só variáveis com prefixo VITE_ viram visíveis, de propósito
// — evita vazar segredo de servidor sem querer). O `?? fallback` cobre o
// caso de alguém rodar sem ter criado o .env ainda.
export const API_BASE_URL = import.meta.env.VITE_API_BASE_URL ?? 'http://172.20.10.2:8080';
export const API_KEY = import.meta.env.VITE_API_KEY ?? 'zica123';
