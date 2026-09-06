import { useEffect, useState } from 'react';

export function useWindowWidth() {
  const [width, setWidth] = useState(1024);

  useEffect(() => {
    const handleResize = () => setWidth(window.innerWidth);
    window.addEventListener('resize', handleResize);
    return () => window.removeEventListener('resize', handleResize);
  }, []);

  return { width, isMobile: width < 768 };
}
