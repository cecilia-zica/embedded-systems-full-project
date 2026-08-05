package monitor.components

import com.raquo.laminar.api.L._

// Modal de confirmação — equivalente ao showDialog<bool>(...) + AlertDialog
// do Flutter e ao ConfirmDialog.jsx da versão React. Diferente do React
// (onde `open` é uma prop booleana e o componente decide "renderizo ou
// não"), aqui é comum a própria tela decidir SE monta esse elemento ou não
// via `child.maybe <-- signal`, então este componente não precisa saber
// nada sobre estar "aberto" — só aparece na árvore quando alguém o inclui.
object ConfirmDialog {
  def apply(
      title: String,
      message: String,
      confirmLabel: String = "Confirmar",
      onCancel: () => Unit,
      onConfirm: () => Unit,
  ): HtmlElement =
    div(
      cls := "modal-backdrop",
      onClick --> (_ => onCancel()),
      div(
        cls := "modal-box",
        // ev.stopPropagation() evita que o clique DENTRO da caixa borbulhe
        // até o backdrop (que fecharia o modal) — mesma ideia do
        // GestureDetector aninhado "engolindo" o toque no Flutter, ou do
        // e.stopPropagation() na versão React.
        onClick --> (ev => ev.stopPropagation()),
        div(cls := "modal-title", title),
        div(cls := "modal-body", message),
        div(
          cls := "modal-actions",
          button(onClick --> (_ => onCancel()), "Cancelar"),
          button(cls := "danger", onClick --> (_ => onConfirm()), confirmLabel),
        ),
      ),
    )
}
