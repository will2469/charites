export function ConstrainedStatefulCarousel({ activeSlide }: any) {
  return <Carousel min={0} max={10} value={activeSlide} className="h-64" />;
}
