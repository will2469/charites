import React, { useState } from 'react';

export function UntransitionedSearchFilter() {
  const [query, setQuery] = useState('');
  const [results, setResults] = useState<string[]>([]);

  return (
    <input
      type="text"
      value={query}
      onChange={(e) => {
        setQuery(e.target.value);
        setResults(expensiveFilter(e.target.value));
      }}
      placeholder="Search items..."
    />
  );
}

function expensiveFilter(q: string): string[] {
  return [q];
}
