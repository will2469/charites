export function ClampedMinmaxGrid() {
  return (
    <div className="grid grid-cols-[repeat(auto-fit,minmax(min(100%,20rem),1fr))] gap-4">
      <div className="card">Kartu Layanan 1</div>
      <div className="card">Kartu Layanan 2</div>
    </div>
  );
}
