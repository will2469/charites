export function ContradictorySlider({ value }: any) {
  return (
    <Slider
      role="slider"
      aria-valuenow={value}
      aria-valuemin={0}
      aria-valuemax={100}
      className="overflow-x-auto snap-x"
    >
      <Slide>Slide A</Slide>
      <Slide>Slide B</Slide>
    </Slider>
  );
}
