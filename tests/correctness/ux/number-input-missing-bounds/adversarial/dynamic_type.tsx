import React from "react";

export function DynamicTypeInput({ dynamicType }: { dynamicType: string }) {
  return <input type={dynamicType} name="quantity" />;
}
