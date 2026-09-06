import React, { useState } from 'react';

export function ModernGridAccordion() {
  const [isOpen, setIsOpen] = useState(false);

  return (
    <div>
      <button onClick={() => setIsOpen(!isOpen)}>Toggle Panel</button>
      <div className={`grid transition-[grid-template-rows] duration-300 ${isOpen ? 'grid-rows-[1fr]' : 'grid-rows-[0fr]'}`}>
        <div className="overflow-hidden">
          <p>Zero-shift modern CSS grid collapsible panel content.</p>
        </div>
      </div>
    </div>
  );
}
