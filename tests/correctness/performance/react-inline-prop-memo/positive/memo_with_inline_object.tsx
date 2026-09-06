import React, { memo } from 'react';

interface UserCardProps {
  user: { id: string; name: string };
  config: { theme: string; compact: boolean };
  onSelect: () => void;
}

const UserCard = memo(({ user, config, onSelect }: UserCardProps) => {
  return <div onClick={onSelect}>{user.name}</div>;
});

export function UserList({ currentUser }: { currentUser: { id: string; name: string } }) {
  return (
    <UserCard
      user={currentUser}
      config={{ theme: 'dark', compact: true }}
      onSelect={() => console.log('selected')}
    />
  );
}
