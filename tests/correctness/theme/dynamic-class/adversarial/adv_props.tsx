// ADV-001: Standard string template without broken utility splicing
export function NormalCard({ className, children }: { className?: string; children: React.ReactNode }) {
  return (
    <div className={`flex flex-col p-6 rounded-xl ${className || ""}`}>
      {children}
    </div>
  );
}
