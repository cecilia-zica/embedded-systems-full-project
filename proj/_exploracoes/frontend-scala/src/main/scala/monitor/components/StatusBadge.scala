package monitor.components

import com.raquo.laminar.api.L._
import monitor.icons.Icons

// Componente reutilizável: transforma um status semântico ('normal' |
// 'alert' | 'error') num "chip" colorido com ícone + texto — mesma ideia do
// StatusBadge.jsx na versão React. Reaproveitado tanto no chip da lista de
// logs quanto em qualquer outro lugar que precise do mesmo tipo de
// indicador, sem duplicar "qual ícone/texto pra cada status".
object StatusBadge {
  private def label(status: String): String = status match {
    case "normal" => "Normal"
    case "alert"  => "Alerta"
    case _        => "Erro Leitura"
  }

  private def icon(status: String): SvgElement = status match {
    case "normal" => Icons.checkCircle2()
    case "alert"  => Icons.alertTriangle()
    case _        => Icons.xCircle()
  }

  def apply(status: String): HtmlElement =
    span(cls := s"chip status-$status", icon(status), label(status))
}
