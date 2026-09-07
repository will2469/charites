export function ModernWrapBreakWord({ apiKey }: { apiKey: string }) {
  return (
    <div className="p-4">
      <div className="whitespace-nowrap wrap-break-word text-sm text-foreground">
        <span>Protected: {apiKey}</span>
      </div>
      <code className="wrap-break-word font-mono text-xs">
        npm install @will2469/charites
      </code>
    </div>
  );
}
