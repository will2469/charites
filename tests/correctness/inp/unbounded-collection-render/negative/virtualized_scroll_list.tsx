import React, { useRef } from 'react';

export function VirtualizedList({ dynamicDataFromApi, rowVirtualizer }: any) {
  const parentRef = useRef<HTMLDivElement>(null);
  return (
    <div ref={parentRef} className="h-96 overflow-y-auto">
      <div style={{ height: `${rowVirtualizer.getTotalSize()}px` }}>
        {rowVirtualizer.getVirtualItems().map((virtualRow: any) => (
          <div key={virtualRow.index}>
            <span>{dynamicDataFromApi[virtualRow.index]?.name}</span>
          </div>
        ))}
      </div>
    </div>
  );
}
