// NEG-001: Semantic chart tokens
export function AccessibleChart() {
  return (
    <div>
      <Bar dataKey="sales" fill="var(--chart-1)" />
      <Line dataKey="visitors" stroke="var(--chart-2)" />
    </div>
  );
}
