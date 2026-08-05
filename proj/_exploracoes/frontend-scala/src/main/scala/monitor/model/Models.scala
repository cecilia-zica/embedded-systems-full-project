package monitor.model

import upickle.default.{ReadWriter, macroRW}

// Formato exato do JSON que o backend Go devolve (ver proj/backend/logging.go
// e controle.go). No proj/web (JS puro) e no api_service.dart, esses dados
// ficam "soltos" — um Map<String, dynamic> no Dart, um objeto qualquer no
// JS. Aqui a gente descreve o formato como case class: se o backend mandar
// um campo faltando ou com tipo errado, o upickle.read() falha com um erro
// claro em vez de deixar `undefined`/`null` vazar silenciosamente pra UI.
//
// "class" é palavra reservada em Scala (usada pra declarar classes), então
// não dá pra nomear o campo `class` como no JSON. A anotação @key diz ao
// upickle "esse campo Scala chamado `classe` corresponde à chave JSON
// `class`" — sem ela, teríamos que usar `` `class` `` com crases toda vez.
case class LogEntry(
    id: Long,
    bpm: Double,
    spo2: Double,
    @upickle.implicits.key("class") classe: Int,
    user_id: String,
    created_at: String,
)
object LogEntry {
  // macroRW gera, em tempo de compilação, o código de leitura E escrita de
  // JSON pra essa case class — o mesmo papel do json.decode/json.encode
  // manual que o Dart faz em cima de um Map, só que aqui é gerado e
  // verificado pelo compilador, não escrito à mão.
  implicit val rw: ReadWriter[LogEntry] = macroRW
}

case class Config(bpm_threshold: Int, alert_enabled: Boolean)
object Config {
  implicit val rw: ReadWriter[Config] = macroRW
}
