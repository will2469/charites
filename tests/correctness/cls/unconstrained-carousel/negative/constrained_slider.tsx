export function ConstrainedSlider({ slides }: any) {
  return (
    <div className="flex overflow-x-auto snap-x h-64 md:h-96 w-full">
      {slides.map((s: any) => (
        <div key={s.id} className="snap-center shrink-0 w-full h-full">
          <img src={s.url} alt={s.title} className="w-full h-full object-cover" />
        </div>
      ))}
    </div>
  );
}

export function SlideAspectCarousel({ slides }: any) {
  return (
    <div className="flex overflow-x-auto snap-x w-full">
      {slides.map((s: any) => (
        <div key={s.id} className="snap-center shrink-0 w-80 aspect-video">
          <img src={s.url} alt={s.title} className="w-full h-full object-cover" />
        </div>
      ))}
    </div>
  );
}
