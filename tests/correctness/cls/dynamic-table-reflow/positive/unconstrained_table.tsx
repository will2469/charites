import React from 'react';

export function UnconstrainedDataTable({ items }: { items: Array<{ id: string; name: string; price: number }> }) {
  return (
    <table className="w-full border">
      <tbody>
        {items.map((it) => (
          <tr key={it.id}>
            <td>{it.name}</td>
            <td>{it.price}</td>
          </tr>
        ))}
      </tbody>
    </table>
  );
}
