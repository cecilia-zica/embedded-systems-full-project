package monitor.components

import com.raquo.laminar.api.L._
import monitor.icons.Icons

// Componente "burro" (sem estado) — equivalente ao LedLegend.jsx da versão
// React. Ícone da lâmpada vem de Icons.scala (Lucide embutido) em vez de
// emoji: fica com aparência idêntica em qualquer sistema operacional. Em
// Laminar não existe uma classe "Component"; qualquer função que devolve um
// HtmlElement já serve — por convenção usamos `apply()` num object pra
// poder chamar como `LedLegend()`.
object LedLegend {
  private case class Item(status: String, pattern: String, meaning: String)

  private val items = List(
    Item("muted", "Aceso fixo", "Esperando leitura"),
    Item("normal", "Pisca devagar", "Normal"),
    Item("alert", "Pisca rápido", "Leitura inconsistente"),
    Item("error", "2 piscadas + pausa", "Alerta — acima do limiar"),
  )

  def apply(): HtmlElement =
    div(
      cls := "card",
      div(
        cls := "card-title",
        Icons.lightbulb(),
        span("O que o LED do sensor significa"),
      ),
      // .map dentro dos filhos de um elemento é o "for" declarativo do
      // Laminar — mesma ideia do `for (final item in items)` do Dart e do
      // `items.map(...)` do React. Como essa lista é fixa (não vem de um
      // Signal), um List[HtmlElement] normal já é aceito como filhos.
      items.map { item =>
        div(
          cls := "legend-row",
          span(cls := s"dot dot-${item.status}"),
          span(b(s"${item.pattern} — "), item.meaning),
        )
      },
    )
}
