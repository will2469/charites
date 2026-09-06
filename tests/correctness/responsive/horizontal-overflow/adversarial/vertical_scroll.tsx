export function VerticalScroll() {
  return (
    <div className="overflow-y-scroll h-64">
      <div className="p-4">
        <span>Vertical scroll is unaffected by horizontal overflow checks</span>
      </div>
    </div>
  );
}
