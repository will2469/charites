import React, { useState } from 'react';

export function WrappedTransitionSearchFilter() {
  const [query, setQuery] = useState('');
  const [results, setResults] = useState<string[]>([]);

  return (
    <input
      type="text"
      value={query}
      onChange={(e) => {
        setQuery(e.target.value);
        React.startTransition(() => {
          setResults(expensiveFilter(e.target.value));
        });
      }}
      placeholder="Search items with transition..."
    />
  );
}

function expensiveFilter(q: string): string[] {
  return [q];
}
