// Adversarial fixture: Synchronous click handler without network mutations
export function CounterButton({ onClick }: { onClick: () => void }) {
  return (
    <button onClick={onClick} className="bg-secondary text-white px-3 py-1">
      Tambah Nilai
    </button>
  );
}
