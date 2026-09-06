export function BoundedFlexChild() {
  return (
    <div className="flex items-center gap-4">
      <div className="min-w-0 w-full">
        <p className="truncate">{userDescription}</p>
      </div>
    </div>
  );
}
