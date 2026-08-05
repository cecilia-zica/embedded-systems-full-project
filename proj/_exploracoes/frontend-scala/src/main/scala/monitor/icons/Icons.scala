package monitor.icons

import com.raquo.laminar.api.L._
import com.raquo.laminar.api.L.{svg => S}

// Ícones SVG da lib Lucide (a mesma usada na versão React, via lucide-react)
// — só que aqui embutidos como código Scala em vez de instalados via npm.
// Scala.js deste projeto não usa nenhum bundler (sem Vite/Webpack, só sbt),
// então não dá pra fazer `import { RefreshCw } from "lucide-react"`; em vez
// disso, copiamos o "d" (o desenho do ícone, em coordenadas SVG) direto do
// pacote lucide-static — que é código aberto (licença ISC) — e desenhamos
// com os elementos <svg>/<path> do próprio Laminar.
//
// Resultado visual idêntico ao da versão React: mesma lib de ícones, dois
// jeitos diferentes de "instalar" ela.
object Icons {

  // Atributos que TODO ícone da Lucide compartilha: contorno (stroke) em vez
  // de preenchimento (fill), 2px de espessura, pontas arredondadas. Extrair
  // isso pra um helper evita repetir essas 6 linhas nas 9 funções abaixo.
  private def iconBase(size: Int, children: SvgElement*): SvgElement =
    S.svg(
      S.viewBox := "0 0 24 24",
      S.width := size.toString,
      S.height := size.toString,
      S.fill := "none",
      S.stroke := "currentColor",
      S.strokeWidth := "2",
      S.strokeLineCap := "round",
      S.strokeLineJoin := "round",
      children.toList,
    )

  private def p(dAttr: String): SvgElement = S.path(S.d := dAttr)

  def refreshCw(size: Int = 18): SvgElement = iconBase(
    size,
    p("M3 12a9 9 0 0 1 9-9 9.75 9.75 0 0 1 6.74 2.74L21 8"),
    p("M21 3v5h-5"),
    p("M21 12a9 9 0 0 1-9 9 9.75 9.75 0 0 1-6.74-2.74L3 16"),
    p("M8 16H3v5"),
  )

  def trash2(size: Int = 18): SvgElement = iconBase(
    size,
    p("M10 11v6"),
    p("M14 11v6"),
    p("M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6"),
    p("M3 6h18"),
    p("M8 6V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"),
  )

  def heartPulse(size: Int = 20): SvgElement = iconBase(
    size,
    p(
      "M2 9.5a5.5 5.5 0 0 1 9.591-3.676.56.56 0 0 0 .818 0A5.49 5.49 0 0 1 22 9.5c0 2.29-1.5 4-3 5.5l-5.492 5.313a2 2 0 0 1-3 .019L5 15c-1.5-1.5-3-3.2-3-5.5"
    ),
    p("M3.22 13H9.5l.5-1 2 4.5 2-7 1.5 3.5h5.27"),
  )

  def lightbulb(size: Int = 18): SvgElement = iconBase(
    size,
    p("M15 14c.2-1 .7-1.7 1.5-2.5 1-.9 1.5-2.2 1.5-3.5A6 6 0 0 0 6 8c0 1 .2 2.2 1.5 3.5.7.7 1.3 1.5 1.5 2.5"),
    p("M9 18h6"),
    p("M10 22h4"),
  )

  def clipboardList(size: Int = 20): SvgElement = iconBase(
    size,
    S.rect(S.x := "8", S.y := "2", S.width := "8", S.height := "4", S.rx := "1", S.ry := "1"),
    p("M16 4h2a2 2 0 0 1 2 2v14a2 2 0 0 1-2 2H6a2 2 0 0 1-2-2V6a2 2 0 0 1 2-2h2"),
    p("M12 11h4"),
    p("M12 16h4"),
    p("M8 11h.01"),
    p("M8 16h.01"),
  )

  def settings(size: Int = 20): SvgElement = iconBase(
    size,
    p(
      "M9.671 4.136a2.34 2.34 0 0 1 4.659 0 2.34 2.34 0 0 0 3.319 1.915 2.34 2.34 0 0 1 2.33 4.033 2.34 2.34 0 0 0 0 3.831 2.34 2.34 0 0 1-2.33 4.033 2.34 2.34 0 0 0-3.319 1.915 2.34 2.34 0 0 1-4.659 0 2.34 2.34 0 0 0-3.32-1.915 2.34 2.34 0 0 1-2.33-4.033 2.34 2.34 0 0 0 0-3.831A2.34 2.34 0 0 1 6.35 6.051a2.34 2.34 0 0 0 3.319-1.915"
    ),
    S.circle(S.cx := "12", S.cy := "12", S.r := "3"),
  )

  def checkCircle2(size: Int = 13): SvgElement = iconBase(
    size,
    S.circle(S.cx := "12", S.cy := "12", S.r := "10"),
    p("m9 12 2 2 4-4"),
  )

  def alertTriangle(size: Int = 13): SvgElement = iconBase(
    size,
    p("m21.73 18-8-14a2 2 0 0 0-3.48 0l-8 14A2 2 0 0 0 4 21h16a2 2 0 0 0 1.73-3"),
    p("M12 9v4"),
    p("M12 17h.01"),
  )

  def xCircle(size: Int = 13): SvgElement = iconBase(
    size,
    S.circle(S.cx := "12", S.cy := "12", S.r := "10"),
    p("m15 9-6 6"),
    p("m9 9 6 6"),
  )

  def wifiOff(size: Int = 14): SvgElement = iconBase(
    size,
    p("M12 20h.01"),
    p("M8.5 16.429a5 5 0 0 1 7 0"),
    p("M5 12.859a10 10 0 0 1 5.17-2.69"),
    p("M19 12.859a10 10 0 0 0-2.007-1.523"),
    p("M2 8.82a15 15 0 0 1 4.177-2.643"),
    p("M22 8.82a15 15 0 0 0-11.288-3.764"),
    p("m2 2 20 20"),
  )
}
