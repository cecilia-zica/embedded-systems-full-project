package monitor.api

import monitor.model.{LogEntry, Config}
import org.scalajs.dom
import org.scalajs.dom.RequestInit
import scala.concurrent.{Future, ExecutionContext}
import scala.scalajs.js
import scala.scalajs.js.Thenable.Implicits._ // dá o .toFuture nas Promises do fetch
import upickle.default._

// Toda a lógica HTTP mora aqui — mesmo papel do services/api_service.dart
// (Dart) e do hooks/*.js na versão React. Diferença central: em Dart/JS os
// métodos devolvem `dynamic`/objeto solto; aqui cada método devolve um tipo
// concreto (Future[List[LogEntry]], por exemplo), então "esqueci de ler um
// campo" vira erro de compilação, não bug descoberto rodando o app.
object Api {

  // quanto tempo esperar antes de desistir de uma requisição. Sem isso, um
  // IP inalcançável (backend desligado, rede errada) deixa o fetch pendurado
  // por mais de 1 minuto antes do navegador desistir sozinho — tempo de
  // sobra pra simplesmente tentar de novo mais tarde, mas ruim pra UI, que
  // fica com spinner infinito em vez de cair pro modo demonstração rápido.
  private val requestTimeoutMs = 4000

  // Future[T] no Scala é o mesmo conceito do Future<T> do Dart e da
  // Promise<T> do JS: "um valor T que ainda não existe, mas vai existir".
  // ExecutionContext implícito = o "em qual fila de tarefas rodar os
  // callbacks" — passamos o padrão global do Scala.js em cada chamador.
  private def request(path: String, method: String, body: Option[String] = None)(implicit
      ec: ExecutionContext
  ): Future[dom.Response] = {
    // AbortController é a forma padrão do navegador de cancelar um fetch em
    // andamento; combinado com setTimeout, vira um "timeout" pro fetch (que
    // não tem essa opção pronta, ao contrário do timeout do http do Dart).
    val controller = new dom.AbortController()
    val timeoutHandle = dom.window.setTimeout(() => controller.abort(), requestTimeoutMs.toDouble)

    val init = new RequestInit {}
    init.method = method.asInstanceOf[dom.HttpMethod]
    init.headers = js.Dictionary("Content-Type" -> "application/json", "X-API-Key" -> Env.apiKey)
    init.signal = controller.signal
    body.foreach(b => init.body = b)

    dom.fetch(s"${Env.apiBaseUrl}$path", init).toFuture
      .flatMap { res =>
        if (res.ok) Future.successful(res)
        else Future.failed(new Exception(s"erro em $path: ${res.status}"))
      }
      .andThen { case _ => dom.window.clearTimeout(timeoutHandle) }
  }

  // GET /api/v1/logging — busca o histórico de leituras
  def getLogs()(implicit ec: ExecutionContext): Future[List[LogEntry]] =
    request("/api/v1/logging", "GET").flatMap(_.text().toFuture).map(read[List[LogEntry]](_))

  // GET /api/v1/controle — busca a configuração atual
  def getConfig()(implicit ec: ExecutionContext): Future[Config] =
    request("/api/v1/controle", "GET").flatMap(_.text().toFuture).map(read[Config](_))

  // POST /api/v1/controle — salva um novo threshold/alerta.
  // write(Config(...)) serializa a case class pra JSON — o mesmo papel do
  // json.encode({...}) no Dart, mas sem montar o Map na mão.
  def postConfig(threshold: Int, alertEnabled: Boolean)(implicit ec: ExecutionContext): Future[Unit] =
    request("/api/v1/controle", "POST", Some(write(Config(threshold, alertEnabled)))).map(_ => ())

  // DELETE /api/v1/logging — apaga todo o histórico de leituras
  def deleteLogs()(implicit ec: ExecutionContext): Future[Unit] =
    request("/api/v1/logging", "DELETE").map(_ => ())
}
