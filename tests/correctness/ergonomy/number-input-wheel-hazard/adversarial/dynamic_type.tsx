import React from "react";

export function DynamicInput({ dynamicType }: { dynamicType: string }) {
  return <input type={dynamicType} name="quantity" />;
}
