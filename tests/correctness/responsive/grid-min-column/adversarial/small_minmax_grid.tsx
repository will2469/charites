export function SmallMinmaxGrid() {
  return (
    <div className="grid grid-cols-[repeat(auto-fit,minmax(150px,1fr))] gap-2">
      <div className="badge">Tag 1</div>
      <div className="badge">Tag 2</div>
    </div>
  );
}
