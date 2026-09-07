export function UnconstrainedStatefulCarousel({ activeSlide }: any) {
  return <Carousel min={0} max={10} value={activeSlide} />;
}
