// Positive fixture: Search input triggers unthrottled fetch directly on keystroke
export function SearchBox() {
  return (
    <div className="relative">
      <input
        type="search"
        placeholder="Cari produk..."
        onChange={(e) => fetchSuggestions(e.target.value)}
        className="w-full px-4 py-2 border rounded"
      />
    </div>
  );
}
