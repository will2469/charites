import React from 'react';

export function DiscreteClickButton() {
  return (
    <button
      onClick={(e) => {
        // Discrete click event mutation on release / click does not fire at 60Hz-120Hz
        e.currentTarget.style.boxShadow = '0 10px 20px rgba(0,0,0,0.2)';
      }}
    >
      Click Me
    </button>
  );
}
