import React from "react";

interface TelemetryPayload {
  ua: string;
  timestamp: number;
}

export function ErrorReporter({ onReport }: { onReport: (data: TelemetryPayload) => void }) {
  const handleReport = () => {
    // Logging user agent for diagnostic telemetry is permitted
    onReport({
      ua: navigator.userAgent,
      timestamp: Date.now(),
    });
  };

  return (
    <button type="button" onClick={handleReport} className="text-xs text-muted-foreground">
      Kirim Laporan Diagnostik
    </button>
  );
}
