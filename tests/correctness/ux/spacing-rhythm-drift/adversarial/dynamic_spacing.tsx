export function DynamicSpacing({ spacing }: { spacing: number }) {
  return (
    <div className={`flex flex-col gap-${spacing}`}>
      <Field />
      <Field />
      <Field />
    </div>
  );
}
