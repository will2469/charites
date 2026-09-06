// Negative fixture: Debounced search callback protects network from keystroke flooding
export function DebouncedSearchBox() {
  const debouncedFetch = useDebouncedCallback((q: string) => {
    fetchSuggestions(q);
  }, 300);

  return (
    <div className="relative">
      <input
        type="search"
        placeholder="Cari produk..."
        onChange={(e) => debouncedFetch(e.target.value)}
        className="w-full px-4 py-2 border rounded"
      />
    </div>
  );
}
