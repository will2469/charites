import { useCallback, useState } from 'react';

export function useUserProfile(userId: string) {
  const [data, setData] = useState(null);

  const refetch = useCallback(() => {
    fetch(`/api/users/${userId}`).then((res) => res.json()).then(setData);
  }, [userId]);

  return { data, refetch };
}
