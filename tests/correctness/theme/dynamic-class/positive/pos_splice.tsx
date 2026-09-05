// POS-001: Dynamic text-${color}-500 utility splicing
export function BrokenButton({ color }: { color: string }) {
  return (
    <button className={`text-${color}-500 font-medium`}>
      Click
    </button>
  );
}
