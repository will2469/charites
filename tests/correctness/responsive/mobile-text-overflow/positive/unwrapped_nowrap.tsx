export function RawApiKeyDisplay({ apiKey }: { apiKey: string }) {
  return (
    <div className="whitespace-nowrap font-mono text-sm text-foreground">
      <span>API Key: {apiKey}</span>
    </div>
  );
}
