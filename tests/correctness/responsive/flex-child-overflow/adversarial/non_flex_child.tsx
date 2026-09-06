export function NonFlexChild() {
  return (
    <div className="grid grid-cols-1 gap-4">
      <div className="w-full">
        <p className="truncate">{longText}</p>
      </div>
    </div>
  );
}
