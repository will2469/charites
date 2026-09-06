import React, { useEffect } from 'react';

export const ResizeWatcher = () => {
  useEffect(() => {
    const handleResize = () => console.log('window resized');
    window.addEventListener('resize', handleResize);
    return () => {
      window.removeEventListener('resize', handleResize);
    };
  }, []);

  return <div>Watching resize cleanly</div>;
};
