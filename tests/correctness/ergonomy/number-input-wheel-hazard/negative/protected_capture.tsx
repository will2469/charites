import React from "react";
import { Input } from "@/components/ui/input";

export function CustomWheelInput({ handleWheel }: { handleWheel: (e: any) => void }) {
  return (
    <Input
      type="number"
      name="amount"
      min="0"
      onWheelCapture={handleWheel}
    />
  );
}
