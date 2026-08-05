// build.sbt é o equivalente ao pubspec.yaml do Flutter ou ao package.json do
// React: declara nome, dependências e como compilar o projeto.

ThisBuild / scalaVersion := "3.3.4"

lazy val root = project
  .in(file("."))
  // Sem esse plugin o sbt compilaria pra bytecode JVM normal (como o
  // proj/backend em Go não, mas como qualquer app Scala de servidor seria).
  // Com ele, todo o código vira JavaScript.
  .enablePlugins(ScalaJSPlugin)
  .settings(
    name := "frontend-scala",

    // Sem isso, o Scala.js compila uma biblioteca genérica, sem ponto de
    // entrada. Com isso, ele chama automaticamente o `main()` do objeto
    // Main assim que o .js é carregado — o equivalente ao
    // `void main() => runApp(...)` do Dart rodando sozinho ao carregar a página.
    scalaJSUseMainModuleInitializer := true,

    libraryDependencies ++= Seq(
      // Laminar: biblioteca de UI reativa pra Scala.js. Não usa JSX nem
      // Virtual DOM como o React — em vez disso, elementos da tela ficam
      // "amarrados" direto a variáveis reativas (Var/Signal), e só o pedaço
      // que mudou é atualizado no DOM real. É outro paradigma de UI, bom
      // pra comparar com os hooks do React.
      "com.raquo" %%% "laminar" % "17.2.0",

      // scalajs-dom: tipagem Scala pras APIs do navegador (fetch, document,
      // etc.) — sem isso, teríamos que chamar tudo via js.Dynamic "cru",
      // sem checagem de tipos nenhuma.
      "org.scala-js" %%% "scalajs-dom" % "2.8.0",

      // upickle: serialização JSON com tipos. Em vez de tratar a resposta
      // da API como um mapa dinâmico (como o proj/web em JS faz com
      // log.bpm, log.spo2 sem nenhuma garantia), a gente descreve o formato
      // com uma case class e o upickle valida/converte pra gente.
      "com.lihaoyi" %%% "upickle" % "3.3.1",
    ),

    // %%% (com 3 %) em vez de %% é o que diz ao sbt "use a versão dessa lib
    // compilada pra Scala.js", não a versão JVM normal.
  )
