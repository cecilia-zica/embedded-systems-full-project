package monitor.components

import com.raquo.laminar.api.L._
import monitor.icons.Icons

// Aviso mostrado quando a tela está usando MockData em vez de dados reais
// (ver LogsScreen.scala / ConfigScreen.scala) — mesma ideia do
// DemoBanner.jsx na versão React. Deixa explícito que os dados na tela são
// fake, bom pra apresentar/testar sem o backend sem se confundir depois.
object DemoBanner {
  def apply(): HtmlElement =
    div(
      cls := "demo-banner",
      Icons.wifiOff(),
      span("Modo demonstração — sem conexão com o backend"),
    )
}
