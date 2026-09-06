// Adversarial fixture: Synchronous boolean state toggle without await
export function ToggleButton({ onToggle }: { onToggle: () => void }) {
  return (
    <button onClick={() => setFlag(true)} className="btn">
      Aktifkan
    </button>
  );
}
