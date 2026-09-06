import React, { useLayoutEffect, useState } from 'react';

export function UserProfile({ userId }: { userId: string }) {
  const [data, setData] = useState<any>(null);

  useLayoutEffect(() => {
    fetch(`/api/users/${userId}`)
      .then((res) => res.json())
      .then((json) => setData(json));
  }, [userId]);

  return <div>{data ? data.name : 'Loading...'}</div>;
}
