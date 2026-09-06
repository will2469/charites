import React from 'react';

export function FixedTable({ items }: { items: Array<{ id: string; name: string }> }) {
  return (
    <table className="w-full table-fixed border">
      <tbody>
        {items.map((it) => (
          <tr key={it.id}>
            <td>{it.name}</td>
          </tr>
        ))}
      </tbody>
    </table>
  );
}

export function ColgroupTable({ items }: { items: Array<{ id: string; name: string }> }) {
  return (
    <table className="w-full border">
      <colgroup>
        <col className="w-full" />
      </colgroup>
      <tbody>
        {items.map((it) => (
          <tr key={it.id}>
            <td>{it.name}</td>
          </tr>
        ))}
      </tbody>
    </table>
  );
}
