package monitor

import com.raquo.laminar.api.L._
import org.scalajs.dom

// Entry point — igual ao `void main() => runApp(const MyApp());` do Flutter
// e ao `ReactDOM.createRoot(...).render(...)` do main.jsx. Como
// `scalaJSUseMainModuleInitializer := true` está no build.sbt, esse
// `main()` roda sozinho assim que o navegador carrega o .js compilado —
// não precisamos chamar ele manualmente em lugar nenhum.
object Main {
  def main(args: Array[String]): Unit = {
    // renderOnDomContentLoaded espera o HTML terminar de carregar antes de
    // montar (equivalente a colocar o <script> no fim do <body>, ou ao
    // DOMContentLoaded do JS puro) — assim `#app` já existe na página
    // quando o Laminar tenta encontrá-lo.
    renderOnDomContentLoaded(dom.document.getElementById("app"), App())
  }
}
