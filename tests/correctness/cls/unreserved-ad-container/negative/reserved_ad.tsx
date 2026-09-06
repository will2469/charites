export function ReservedAdSlots() {
  return (
    <div className="container mx-auto">
      <div id="ad-leaderboard" data-ad-slot="12345" className="w-full min-h-[90px] md:min-h-[250px] bg-muted/20" />
      <div id="ad-sidebar" data-ad-slot="67890" className="w-full">
        <div className="skeleton animate-pulse w-full h-[250px]" />
      </div>
    </div>
  );
}
