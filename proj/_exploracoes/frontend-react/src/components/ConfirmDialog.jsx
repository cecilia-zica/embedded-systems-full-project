// Modal de confirmação genérico — equivalente ao showDialog<bool>(...) +
// AlertDialog do Flutter. `open` controla se aparece; onConfirm/onCancel
// são callbacks que quem usa o componente decide o que fazer — o modal em
// si não sabe nada sobre "limpar logs" especificamente, só sabe mostrar
// título/mensagem/2 botões. Isso é o que deixa ele reaproveitável.
export default function ConfirmDialog({ open, title, message, confirmLabel = 'Confirmar', onConfirm, onCancel }) {
  if (!open) return null;

  return (
    <div className="modal-backdrop" onClick={onCancel}>
      {/* stopPropagation: sem isso, clicar DENTRO da caixa também contaria
          como clique no backdrop (que fecha o modal). */}
      <div className="modal-box" onClick={(e) => e.stopPropagation()}>
        <div className="modal-title">{title}</div>
        <div className="modal-body">{message}</div>
        <div className="modal-actions">
          <button onClick={onCancel}>Cancelar</button>
          <button className="danger" onClick={onConfirm}>
            {confirmLabel}
          </button>
        </div>
      </div>
    </div>
  );
}
