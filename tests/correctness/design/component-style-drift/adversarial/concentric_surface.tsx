import React from "react";

// A2: Concentric radius invariant - Card legitimately uses rounded-2xl without conflicting with Button's rounded-md
export function ConcentricSurface() {
  return (
    <Card className="rounded-2xl bg-black text-white p-8">
      <div className="flex justify-end">
        <Button className="rounded-md bg-white text-black px-4 py-2">
          Inner Action
        </Button>
      </div>
    </Card>
  );
}
