export function DesktopOnlyMarketingBadge() {
  return (
    <div className="flex items-center gap-2">
      <span>Fitur Desa Digital</span>
      <span className="hidden md:inline-flex px-2 py-1 bg-muted text-xs rounded">
        Versi Desktop
      </span>
    </div>
  );
}
