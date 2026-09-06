import React, { useLayoutEffect, useRef, useState } from 'react';

export function TooltipBubble() {
  const ref = useRef<HTMLDivElement>(null);
  const [coords, setCoords] = useState({ top: 0, left: 0 });

  useLayoutEffect(() => {
    if (ref.current) {
      const rect = ref.current.getBoundingClientRect();
      setCoords({ top: rect.top, left: rect.left });
    }
  }, []);

  return <div ref={ref} style={{ top: coords.top, left: coords.left }}>Tooltip Content</div>;
}
