// sbt-scalajs = o plugin que ensina o sbt a compilar Scala pra JavaScript em
// vez de bytecode JVM. Sem isso, `sbt compile` geraria .class normais — com
// ele, `sbt fastLinkJS` gera um arquivo .js que o navegador consegue rodar.
addSbtPlugin("org.scala-js" % "sbt-scalajs" % "1.17.0")
