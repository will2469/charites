import React from 'react';

export function StaticDocTable() {
  return (
    <table className="w-full border">
      <thead>
        <tr>
          <th>Param</th>
          <th>Type</th>
          <th>Description</th>
        </tr>
      </thead>
      <tbody>
        <tr>
          <td>id</td>
          <td>string</td>
          <td>Unique record identifier</td>
        </tr>
        <tr>
          <td>status</td>
          <td>boolean</td>
          <td>Active flag</td>
        </tr>
      </tbody>
    </table>
  );
}
