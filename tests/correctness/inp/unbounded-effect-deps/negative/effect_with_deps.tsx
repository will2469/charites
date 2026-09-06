import React, { useEffect } from 'react';

export function BoundedEffectComponent() {
  useEffect(() => {
    recomputeLayoutOnEveryRender();
  }, []);

  return <div>Component Content</div>;
}

function recomputeLayoutOnEveryRender() {
  // layout
}
