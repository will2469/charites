export function ScrollableToolbar() {
  return (
    <div className="flex items-center gap-2 p-2 bg-surface overflow-x-auto">
      <button type="button">Edit</button>
      <button type="button">Salin</button>
      <button type="button">Cetak</button>
      <button type="button">Bagikan</button>
      <button type="button">Hapus</button>
    </div>
  );
}
