import React from 'react';

function StandardCard({ data }: { data: any }) {
  return <div>{data.title}</div>;
}

export function App() {
  return (
    <div>
      <button onClick={() => console.log('clicked')}>Click</button>
      <StandardCard data={{ title: 'Non-memoized' }} />
    </div>
  );
}
