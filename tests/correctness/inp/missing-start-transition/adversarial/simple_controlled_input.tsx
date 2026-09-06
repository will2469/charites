import React, { useState } from 'react';

export function SimpleControlledInput() {
  const [query, setQuery] = useState('');

  return (
    <input
      type="text"
      value={query}
      onChange={(e) => setQuery(e.target.value)}
      placeholder="Simple input..."
    />
  );
}
