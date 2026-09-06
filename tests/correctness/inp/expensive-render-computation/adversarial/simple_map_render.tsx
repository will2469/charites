import React from 'react';

interface User {
  id: string;
  name: string;
}

export function SimpleUserList({ users }: { users: User[] }) {
  const userNames = users.map((u) => u.name);

  return (
    <ul>
      {userNames.map((name, idx) => (
        <li key={idx}>{name}</li>
      ))}
    </ul>
  );
}
