export function SizedHeroImage() {
  return (
    <div className="container mx-auto">
      <img src="/hero.jpg" alt="Hero Banner" width={1200} height={600} className="w-full h-auto" />
      <img src="/avatar.jpg" alt="Avatar" className="w-10 h-10 rounded-full" />
      <img src="/thumb.jpg" alt="Thumbnail" className="w-full aspect-video object-cover" />
    </div>
  );
}
