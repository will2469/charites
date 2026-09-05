// NEG-001: Static lookup dictionary
const colorMap: Record<string, string> = {
  red: "text-red-500",
  blue: "text-blue-500",
};

export function SafeButton({ color }: { color: string }) {
  return (
    <button className={`${colorMap[color]} font-medium`}>
      Click
    </button>
  );
}
