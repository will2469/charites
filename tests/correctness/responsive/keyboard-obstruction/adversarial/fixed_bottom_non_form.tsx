export function FixedBottomNonForm() {
  return (
    <div className="fixed bottom-0 inset-x-0 bg-surface border-t p-4 flex justify-between items-center z-50">
      <p className="text-sm">Kami menggunakan cookie untuk meningkatkan pengalaman warga.</p>
      <button type="button" className="px-4 py-2 bg-primary text-white rounded">
        Setuju
      </button>
    </div>
  );
}
