package monitor

import com.raquo.laminar.api.L._
import monitor.screens.{LogsScreen, ConfigScreen}
import monitor.icons.Icons

// Componente raiz — equivalente a MyApp + HomeNav + _HomeNavState no Dart e
// a App.jsx na versão React: guarda qual aba está ativa e desenha a
// bottom nav.
object App {
  def apply(): HtmlElement = {
    // Var(0) = o `int _indiceAtual = 0;` do Dart / `useState(0)` do React —
    // guarda qual aba está selecionada.
    val tabVar = Var(0)

    div(
      styleAttr := "display:flex;flex-direction:column;min-height:100vh",
      div(
        styleAttr := "flex:1",
        // Troca a tela inteira conforme a aba — mesma ideia do
        // `_telas[_indiceAtual]` do Dart e do `tab === 0 ? <LogsScreen/> : <ConfigScreen/>` do React.
        child <-- tabVar.signal.map(tab => if (tab == 0) LogsScreen() else ConfigScreen()),
      ),
      div(
        cls := "bottom-nav",
        button(
          cls <-- tabVar.signal.map(tab => s"nav-item ${if (tab == 0) "active" else ""}"),
          onClick --> (_ => tabVar.set(0)),
          Icons.clipboardList(),
          span("Logs"),
        ),
        button(
          cls <-- tabVar.signal.map(tab => s"nav-item ${if (tab == 1) "active" else ""}"),
          onClick --> (_ => tabVar.set(1)),
          Icons.settings(),
          span("Config"),
        ),
      ),
    )
  }
}
