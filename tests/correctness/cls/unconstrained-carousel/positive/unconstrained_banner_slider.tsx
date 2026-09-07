export function UnconstrainedBannerSlider({ onSlideChange }: any) {
  return (
    <BannerSlider onChange={onSlideChange}>
      <div>Slide 1</div>
      <div>Slide 2</div>
    </BannerSlider>
  );
}
