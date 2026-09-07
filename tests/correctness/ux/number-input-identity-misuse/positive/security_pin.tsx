import React from "react";
import { Input } from "@/components/ui/input";

export function AuthSecurityPin() {
  return (
    <div>
      <Input type="number" name="security_pin" />
    </div>
  );
}
