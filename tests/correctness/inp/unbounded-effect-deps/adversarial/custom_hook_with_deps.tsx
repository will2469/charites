import React from 'react';

export function CustomHookComponent() {
  useCustomLogger(() => {
    logRender();
  });

  return <div>Component Content</div>;
}

function useCustomLogger(fn: () => void) {
  fn();
}

function logRender() {
  // log
}
