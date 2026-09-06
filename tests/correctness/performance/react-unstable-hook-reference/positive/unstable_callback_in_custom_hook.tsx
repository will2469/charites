import { useState } from 'react';

export function useUserProfile(userId: string) {
  const [data, setData] = useState(null);

  const refetch = () => {
    fetch(`/api/users/${userId}`).then((res) => res.json()).then(setData);
  };

  return { data, refetch };
}
