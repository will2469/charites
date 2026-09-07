export function AccessibleCustomSlider({ value }: any) {
  return (
    <div
      role="slider"
      aria-valuenow={value}
      aria-valuemin={0}
      aria-valuemax={100}
      tabIndex={0}
    />
  );
}
