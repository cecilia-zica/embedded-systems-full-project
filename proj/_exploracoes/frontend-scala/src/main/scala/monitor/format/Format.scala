package monitor.format

import scala.scalajs.js
import js.Date

// Funções puras, sem UI e sem estado — mesma lógica que _classToText,
// _classToColor e _formatTimestamp no Dart. Diferença pra versão anterior:
// classToColor devolvia um hex cru; classToStatus devolve uma palavra
// semântica ('normal' | 'alert' | 'error'). Quem decide a cor de verdade é
// o CSS (ver styles.css, classes .status-*) — separa "o que o dado
// significa" de "como ele é desenhado", igual foi feito na versão React.
object Format {

  // 0 = Repouso/Normal, 1 = Ativo/Alerta, 2 = Erro — mesma classificação do
  // TinyML rodando no ESP32 (proj/firmware/src/classificacao.h)
  def classToStatus(c: Int): String = c match {
    case 0 => "normal"
    case 1 => "alert"
    case _ => "error"
  }

  // js.Date é a classe Date do JavaScript "importada" pro Scala.js — igual
  // o `new Date(iso)` do JS, entende o "Z" (UTC) do timestamp do backend e
  // os getters (getHours, getDate...) já devolvem hora local do navegador
  // sozinhos (mesmo resultado do `DateTime.parse(iso).toLocal()` do Dart).
  def formatTimestamp(iso: String): String = {
    val dt = new Date(iso)
    def two(n: Double): String = f"${n.toInt}%02d"
    s"${two(dt.getDate())}/${two(dt.getMonth() + 1)} ${two(dt.getHours())}:${two(dt.getMinutes())}:${two(dt.getSeconds())}"
  }
}
