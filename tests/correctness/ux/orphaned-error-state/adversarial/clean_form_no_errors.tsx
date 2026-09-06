// Adversarial fixture: Form without client-side error setters
export function SimpleSearchForm() {
  return (
    <form action="/search" method="GET" className="flex gap-2">
      <input type="search" name="q" placeholder="Cari..." className="border px-3 py-1" />
      <button type="submit" className="bg-primary text-white px-3 py-1">Cari</button>
    </form>
  );
}
