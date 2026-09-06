export function DesktopNowrapTable() {
  return (
    <div className="flex flex-col gap-2">
      <span className="md:whitespace-nowrap text-sm text-foreground">
        Kolom ini hanya nowrap pada layar desktop (md ke atas).
      </span>
    </div>
  );
}
