// POS-001: Hardcoded fill and stroke in Recharts components
export function SalesChart() {
  return (
    <div>
      <Bar dataKey="sales" fill="#3b82f6" />
      <Line dataKey="visitors" stroke="#10b981" />
    </div>
  );
}
