import React from 'react';

export function SlicedList({ dynamicDataFromApi }: { dynamicDataFromApi: Array<{ id: string; name: string }> }) {
  return (
    <div className="h-96 overflow-y-auto">
      {dynamicDataFromApi.slice(0, 20).map(item => (
        <div key={item.id}>
          <span>{item.name}</span>
        </div>
      ))}
    </div>
  );
}
