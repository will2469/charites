// Adversarial fixture: Try block with only finally and no catch
export function ActionButton({ onClean }: { onClean: () => void }) {
  return (
    <button
      onClick={() => {
        try {
          doCleanup();
        } finally {
          releaseLock();
        }
      }}
      className="btn"
    >
      Bersihkan
    </button>
  );
}
