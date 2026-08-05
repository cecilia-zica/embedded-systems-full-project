package monitor.api

import monitor.model.{LogEntry, Config}
import scala.scalajs.js
import scala.scalajs.js.Date

// Dados de demonstração — usados quando o fetch pro backend falha (ver
// LogsScreen.scala e ConfigScreen.scala). Existe pra você conseguir ver e
// mexer no app mesmo sem o backend Go rodando ou sem o ESP32 conectado —
// mesma ideia do mockData.js na versão React, linha por linha.
object MockData {

  // Timestamps relativos a "agora" (Date.now() - X minutos) em vez de datas
  // fixas: assim os logs sempre parecem recentes, não importa quando você
  // abrir o app.
  private def minutesAgo(minutes: Int): String = {
    val ms = js.Date.now() - minutes * 60000.0
    new Date(ms).toISOString()
  }

  val logs: List[LogEntry] = List(
    LogEntry(8, 76, 98, 0, "A1B2C3D4", minutesAgo(2)),
    LogEntry(7, 132, 95, 1, "A1B2C3D4", minutesAgo(9)),
    LogEntry(6, 81, 97, 0, "F09E1122", minutesAgo(18)),
    LogEntry(5, 0, 0, 2, "7C3A9B00", minutesAgo(25)),
    LogEntry(4, 74, 98, 0, "A1B2C3D4", minutesAgo(40)),
    LogEntry(3, 145, 93, 1, "D41D8CD9", minutesAgo(55)),
    LogEntry(2, 79, 99, 0, "F09E1122", minutesAgo(70)),
    LogEntry(1, 88, 97, 0, "A1B2C3D4", minutesAgo(95)),
  )

  val config: Config = Config(bpm_threshold = 120, alert_enabled = true)
}
