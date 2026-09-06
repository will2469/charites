export function ResponsiveGrid() {
  return (
    <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-4 gap-4">
      <div className="bg-card p-4">Item 1</div>
      <div className="bg-card p-4">Item 2</div>
    </div>
  );
}
