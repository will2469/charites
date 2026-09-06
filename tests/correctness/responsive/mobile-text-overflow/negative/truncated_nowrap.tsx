export function TruncatedApiKeyDisplay({ apiKey }: { apiKey: string }) {
  return (
    <div className="whitespace-nowrap truncate font-mono text-sm text-foreground">
      <span>API Key: {apiKey}</span>
    </div>
  );
}
