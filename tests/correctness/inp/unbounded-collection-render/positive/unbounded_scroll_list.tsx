import React from 'react';

export function UnboundedList({ dynamicDataFromApi, handleSelect }: { dynamicDataFromApi: Array<{ id: string; name: string }>; handleSelect: (id: string) => void }) {
  return (
    <div className="h-96 overflow-y-auto">
      {dynamicDataFromApi.map(item => (
        <div key={item.id} onClick={() => handleSelect(item.id)}>
          <span>{item.name}</span>
        </div>
      ))}
    </div>
  );
}
