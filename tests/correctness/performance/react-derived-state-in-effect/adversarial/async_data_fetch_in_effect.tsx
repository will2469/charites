import React, { useEffect, useState } from 'react';

export const UserDetailsCard = ({ userId }: { userId: string }) => {
  const [user, setUser] = useState<any>(null);

  useEffect(() => {
    fetch(`/api/users/${userId}`)
      .then((res) => res.json())
      .then((data) => setUser(data));
  }, [userId]);

  return <div>User: {user ? user.name : 'Loading...'}</div>;
};
