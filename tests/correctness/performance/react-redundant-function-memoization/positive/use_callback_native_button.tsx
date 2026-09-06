import React, { useCallback, useState } from 'react';

export const SimpleCounter = () => {
  const [count, setCount] = useState(0);

  const handleIncrement = useCallback(() => {
    setCount((c) => c + 1);
  }, []);

  return (
    <div>
      <span>Count: {count}</span>
      <button onClick={handleIncrement}>Increment</button>
    </div>
  );
};
