import React, { memo, useCallback } from 'react';

interface UserCardProps {
  user: { id: string; name: string };
  config: { theme: string; compact: boolean };
  onSelect: () => void;
}

const USER_CONFIG = { theme: 'dark', compact: true } as const;

const UserCard = memo(({ user, config, onSelect }: UserCardProps) => {
  return <div onClick={onSelect}>{user.name}</div>;
});

export function UserList({ currentUser }: { currentUser: { id: string; name: string } }) {
  const handleSelect = useCallback(() => {
    console.log('selected');
  }, []);

  return (
    <UserCard
      user={currentUser}
      config={USER_CONFIG}
      onSelect={handleSelect}
    />
  );
}
