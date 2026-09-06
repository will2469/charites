export function AutoScrollWrapped() {
  return (
    <div className="w-full overflow-x-auto">
      <div className="flex gap-4 min-w-max">
        <div className="p-4 bg-card">Item 1</div>
        <div className="p-4 bg-card">Item 2</div>
      </div>
    </div>
  );
}
