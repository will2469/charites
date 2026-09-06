import React, { useState } from 'react';

export function LegacyAccordion() {
  const [isOpen, setIsOpen] = useState(false);

  return (
    <div>
      <button onClick={() => setIsOpen(!isOpen)}>Toggle Panel</button>
      <div className={`transition-all duration-300 overflow-hidden ${isOpen ? 'max-h-[1000px]' : 'max-h-0'}`}>
        <p>Panel content inside arbitrarily large max-height container.</p>
      </div>
    </div>
  );
}
