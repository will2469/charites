export function VolumeSettings({ volume, setVolume }: any) {
  return (
    <Slider
      min={0}
      max={100}
      step={1}
      value={volume}
      onValueChange={setVolume}
    />
  );
}
