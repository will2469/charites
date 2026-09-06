import React from "react";

// A2: Plain properties or variables with similar names do not invoke experimental Web APIs
export function SocialStats({ shareCount, hasBluetoothDevice }: { shareCount: number; hasBluetoothDevice: boolean }) {
  return (
    <div>
      <span>Total Shares: {shareCount}</span>
      <span>Connected Device: {hasBluetoothDevice ? "Yes" : "No"}</span>
    </div>
  );
}
