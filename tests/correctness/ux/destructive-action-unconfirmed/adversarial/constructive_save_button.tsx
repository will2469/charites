// Adversarial fixture: Standard non-destructive save button
export function SaveButton({ onSave }: { onSave: () => void }) {
  return (
    <button onClick={onSave} className="bg-primary text-white px-4 py-2">
      Simpan Data
    </button>
  );
}
