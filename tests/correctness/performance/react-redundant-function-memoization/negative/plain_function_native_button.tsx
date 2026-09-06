import React, { useState } from 'react';

export const SimpleCounter = () => {
  const [count, setCount] = useState(0);

  const handleIncrement = () => {
    setCount((c) => c + 1);
  };

  return (
    <div>
      <span>Count: {count}</span>
      <button onClick={handleIncrement}>Increment</button>
    </div>
  );
};
