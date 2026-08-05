package monitor.screens

import com.raquo.laminar.api.L._
import monitor.api.{Api, MockData}
import monitor.model.LogEntry
import monitor.format.Format
import monitor.icons.Icons
import monitor.components.{LedLegend, ConfirmDialog, StatusBadge, DemoBanner}
import scala.concurrent.ExecutionContext.Implicits.global
import scala.util.{Success, Failure}

// Tela 1 — mostra o histórico (logs) de leituras que o ESP32 mandou.
// Equivalente a LogsScreen + _LogsScreenState no Dart e a LogsScreen.jsx na
// versão React.
object LogsScreen {
  def apply(): HtmlElement = {
    // Var = "caixa" de estado mutável e reativa — o equivalente Laminar do
    // useState do React / das variáveis dentro de _LogsScreenState no Dart.
    // Diferença importante: no React, mudar o estado redesenha TODO o
    // componente de novo (função inteira reexecuta); em Laminar, só os
    // elementos amarrados via `<--` àquela Var específica são atualizados
    // no DOM real — o resto da árvore nem "sabe" que algo mudou.
    val logsVar = Var(List.empty[LogEntry])
    val loadingVar = Var(true)
    val erroVar = Var(Option.empty[String])
    val confirmandoLimparVar = Var(false)
    // true quando o backend não respondeu e a tela está mostrando
    // MockData.logs em vez de dados reais — a UI usa isso pra exibir o
    // aviso de "modo demonstração" (DemoBanner) em vez de fingir que são
    // dados reais.
    val usingMockVar = Var(false)

    def fetchLogs(): Unit =
      Api.getLogs().onComplete {
        case Success(data) =>
          logsVar.set(data)
          loadingVar.set(false)
          erroVar.set(None)
          usingMockVar.set(false)
        case Failure(_) =>
          // Backend fora do ar (ou IP errado, ou você só quer ver o app sem
          // ligar nada) — em vez de deixar a tela com uma mensagem de erro,
          // cai pros dados de demonstração (ver MockData.scala).
          logsVar.set(MockData.logs)
          loadingVar.set(false)
          erroVar.set(None)
          usingMockVar.set(true)
      }

    def confirmarLimpar(): Unit = {
      confirmandoLimparVar.set(false)
      // em modo demonstração não existe backend pra chamar — só limpa local.
      if (usingMockVar.now()) {
        logsVar.set(Nil)
      } else {
        Api.deleteLogs().onComplete {
          case Success(_)   => logsVar.set(Nil)
          case Failure(err) => erroVar.set(Some(err.getMessage))
        }
      }
    }

    def logRow(log: LogEntry): HtmlElement = {
      val status = Format.classToStatus(log.classe)
      div(
        cls := "log-item",
        div(cls := s"avatar status-$status", Icons.heartPulse()),
        div(
          cls := "log-main",
          div(cls := "log-title", s"${log.bpm.toInt} BPM · SpO2 ${log.spo2.toInt}%"),
          div(cls := "log-subtitle", s"${log.user_id} — ${Format.formatTimestamp(log.created_at)}"),
        ),
        StatusBadge(status),
      )
    }

    // combineWith junta vários Signals num Signal só de tupla — precisamos
    // dos 3 estados (loading/erro/logs) ao mesmo tempo pra decidir o que
    // desenhar, igual o if/else if encadeado que lê _loading/_erro/_logs.length
    // no Dart e no React (só que lá são 3 variáveis soltas, aqui é 1 Signal combinado).
    val bodySignal = loadingVar.signal.combineWith(erroVar.signal, logsVar.signal).map {
      case (loading, erro, logs) =>
        if (loading) div(cls := "spinner")
        else if (erro.isDefined) div(cls := "center-msg", s"Erro: ${erro.get}")
        else if (logs.isEmpty) div(cls := "center-msg", "Nenhuma leitura ainda")
        else div(logs.map(logRow))
    }

    div(
      styleAttr := "display:flex;flex-direction:column;height:100%",
      // onMountCallback roda quando este elemento entra de fato no DOM —
      // é o equivalente ao initState() do Flutter e ao
      // useEffect(fn, []) do React: "roda uma vez, quando a tela nasce".
      onMountCallback(_ => fetchLogs()),
      div(
        cls := "appbar",
        "Logs",
        div(
          cls := "appbar-actions",
          button(cls := "icon-btn", title := "Atualizar", onClick --> (_ => fetchLogs()), Icons.refreshCw()),
          button(
            cls := "icon-btn",
            title := "Limpar histórico",
            onClick --> (_ => confirmandoLimparVar.set(true)),
            Icons.trash2(),
          ),
        ),
      ),
      div(
        cls := "content",
        child.maybe <-- usingMockVar.signal.map(if (_) Some(DemoBanner()) else None),
        LedLegend(),
        // `child <-- signal` amarra ESTE ponto da árvore ao valor mais
        // recente do Signal: toda vez que bodySignal emite um novo
        // HtmlElement, o Laminar troca só esse pedaço no DOM real.
        child <-- bodySignal,
      ),
      // `child.maybe <-- signalDeOption` só desenha algo quando o Option é
      // Some — mesma ideia do `snackbar && <div>...` do React ou do
      // `if (confirmandoLimpar) ...` do showDialog no Flutter.
      child.maybe <-- confirmandoLimparVar.signal.map { open =>
        if (open)
          Some(
            ConfirmDialog(
              title = "Limpar log?",
              message = "Apaga todo o histórico de leituras. Não dá pra desfazer.",
              confirmLabel = "Limpar",
              onCancel = () => confirmandoLimparVar.set(false),
              onConfirm = () => confirmarLimpar(),
            )
          )
        else None
      },
    )
  }
}
