import React from "react";
import { Input } from "@/components/ui/input";

export function DonationForm() {
  return (
    <div>
      <Input type="number" name="amount" placeholder="Enter amount" />
    </div>
  );
}
