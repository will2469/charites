export function UnconstrainedSlider({ slides }: any) {
  return (
    <div className="flex overflow-x-auto snap-x">
      {slides.map((s: any) => <img key={s.id} src={s.url} alt={s.title} />)}
    </div>
  );
}
