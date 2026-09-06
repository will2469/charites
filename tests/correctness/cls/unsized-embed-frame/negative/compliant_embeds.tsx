export function CompliantEmbeds() {
  return (
    <div className="container mx-auto">
      <div className="w-full aspect-video">
        <iframe src="https://www.youtube.com/embed/dQw4w9WgXcQ" title="Video" className="w-full h-full" />
      </div>
      <video src="/desa-profile.mp4" width={640} height={360} controls className="w-full h-auto" />
      <div className="min-h-[300px]">
        <iframe src="/map" title="Peta Desa" className="w-full h-full" />
      </div>
    </div>
  );
}
