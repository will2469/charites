import React, { useEffect } from 'react';

export function UnboundedEffectComponent() {
  useEffect(() => {
    recomputeLayoutOnEveryRender();
  });

  return <div>Component Content</div>;
}

function recomputeLayoutOnEveryRender() {
  // layout
}
