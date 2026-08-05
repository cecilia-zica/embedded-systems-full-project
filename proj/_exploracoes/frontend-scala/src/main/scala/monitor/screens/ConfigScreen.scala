package monitor.screens

import com.raquo.laminar.api.L._
import monitor.api.{Api, MockData}
import monitor.components.DemoBanner
import org.scalajs.dom
import scala.concurrent.ExecutionContext.Implicits.global
import scala.util.{Success, Failure}

// Tela 2 — mostra e edita o limiar de alerta (bpm_threshold) que o ESP32
// consulta a cada 30s. Equivalente a ConfigScreen + _ConfigScreenState no
// Dart e a ConfigScreen.jsx na versão React.
object ConfigScreen {
  def apply(): HtmlElement = {
    val thresholdVar = Var("")
    val alertEnabledVar = Var(true)
    val loadingVar = Var(true)
    val savingVar = Var(false)
    val snackbarVar = Var(Option.empty[String]) // mensagem temporária, tipo o SnackBar do Flutter
    val usingMockVar = Var(false)

    def showSnackbar(msg: String): Unit = {
      snackbarVar.set(Some(msg))
      // window.setTimeout é o mesmo setTimeout do JS; some sozinho depois de
      // 3s, igual o comportamento padrão do SnackBar do Flutter.
      dom.window.setTimeout(() => snackbarVar.set(None), 3000)
    }

    def loadCurrentConfig(): Unit =
      Api.getConfig().onComplete {
        case Success(config) =>
          thresholdVar.set(config.bpm_threshold.toString)
          alertEnabledVar.set(config.alert_enabled)
          loadingVar.set(false)
          usingMockVar.set(false)
        case Failure(_) =>
          // backend fora do ar — pré-preenche com a config de demonstração
          // em vez de deixar o formulário vazio (ver MockData.scala).
          thresholdVar.set(MockData.config.bpm_threshold.toString)
          alertEnabledVar.set(MockData.config.alert_enabled)
          loadingVar.set(false)
          usingMockVar.set(true)
      }

    def saveConfig(): Unit =
      // toIntOption é o tryParse do Dart / o Number.parseInt-que-não-lança
      // do React: devolve None em vez de estourar exceção quando o texto
      // não é um número válido — valida ANTES de gastar uma chamada de rede.
      thresholdVar.now().toIntOption match {
        case None =>
          showSnackbar("Digite um número válido pro limiar de BPM")
        case Some(parsed) =>
          savingVar.set(true)
          // sem backend de verdade pra salvar em modo demonstração — só
          // simula um delay de rede, pra dar pra testar a interação mesmo offline.
          if (usingMockVar.now()) {
            dom.window.setTimeout(
              () => {
                savingVar.set(false)
                showSnackbar("Configuração salva! (modo demonstração)")
              },
              400,
            )
          } else {
            Api.postConfig(parsed, alertEnabledVar.now()).onComplete {
              case Success(_) =>
                savingVar.set(false)
                showSnackbar("Configuração salva!")
              case Failure(err) =>
                savingVar.set(false)
                showSnackbar(s"Erro ao salvar: ${err.getMessage}")
            }
          }
      }

    div(
      onMountCallback(_ => loadCurrentConfig()),
      child <-- loadingVar.signal.map { loading =>
        if (loading) div(cls := "spinner")
        else
          div(
            div(cls := "appbar", "Configurações"),
            child.maybe <-- usingMockVar.signal.map(if (_) Some(DemoBanner()) else None),
            div(
              cls := "card",
              label(cls := "field-label", "Limiar de alerta (BPM)"),
              input(
                cls := "text-input",
                typ := "number",
                placeholder := "ex: 120",
                // value <-- mantém o <input> sincronizado com a Var (fonte
                // única de verdade); onInput.mapToValue --> escreve de volta
                // na Var a cada tecla. Junto, isso é um "input controlado" —
                // mesma ideia do TextEditingController do Dart e do
                // useState+onChange do React, só que expresso como um fluxo
                // de dados nos dois sentidos em vez de "estado + callback".
                value <-- thresholdVar.signal,
                onInput.mapToValue --> thresholdVar,
              ),
              div(
                cls := "switch-row",
                span("Alerta ativado"),
                button(
                  cls <-- alertEnabledVar.signal.map(on => s"switch ${if (on) "on" else ""}"),
                  onClick --> (_ => alertEnabledVar.update(!_)),
                  span(cls := "switch-thumb"),
                ),
              ),
            ),
            button(
              cls := "btn-primary",
              disabled <-- savingVar.signal,
              onClick --> (_ => saveConfig()),
              child.text <-- savingVar.signal.map(s => if (s) "Salvando..." else "SALVAR"),
            ),
            child.maybe <-- snackbarVar.signal.map(_.map(msg => div(cls := "snackbar", msg))),
          )
      },
    )
  }
}
