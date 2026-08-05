package monitor.api

import scala.scalajs.js

// Fonte única de configuração do app — mesmo papel do src/config.js na
// versão React. Scala.js não tem um sistema de ".env" pronto (isso é coisa
// de bundler tipo Vite/Webpack, e aqui a gente compila só com sbt, sem
// bundler); o jeito equivalente é o index.html definir `window.__CONFIG__`
// ANTES de carregar o main.js (ver a tag <script> em index.html), e este
// objeto só lê esse valor. Assim nenhum outro arquivo Scala precisa saber
// que a configuração "vem de fora" — só este.
object Env {
  // js.Dynamic.global.window é o `window` do navegador, sem tipagem — a
  // única forma de acessar algo que não faz parte da API padrão do DOM
  // (scalajs-dom não tem como saber de um `__CONFIG__` que a gente inventou).
  private val raw = js.Dynamic.global.window.__CONFIG__

  // `?? fallback`: cobre o caso de abrir o app sem passar pelo index.html
  // certo (ex: durante testes) — mesma ideia do `?? fallback` no config.js do React.
  val apiBaseUrl: String =
    if (js.isUndefined(raw)) "http://172.20.10.2:8080" else raw.apiBaseUrl.asInstanceOf[String]

  val apiKey: String =
    if (js.isUndefined(raw)) "zica123" else raw.apiKey.asInstanceOf[String]
}
