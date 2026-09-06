import React, { useMemo } from 'react';

interface User {
  id: string;
  name: string;
  score: number;
}

export function UserList({ users, filterText }: { users: User[]; filterText: string }) {
  const visibleUsers = useMemo(() => {
    return users
      .filter((u) => u.name.includes(filterText))
      .sort((a, b) => b.score - a.score);
  }, [users, filterText]);

  return (
    <ul>
      {visibleUsers.map((u) => (
        <li key={u.id}>{u.name}</li>
      ))}
    </ul>
  );
}
