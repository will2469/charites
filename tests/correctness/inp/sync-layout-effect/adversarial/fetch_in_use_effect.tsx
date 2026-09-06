import React, { useEffect, useState } from 'react';

export function UserProfileAsync({ userId }: { userId: string }) {
  const [data, setData] = useState<any>(null);

  useEffect(() => {
    fetch(`/api/users/${userId}`)
      .then((res) => res.json())
      .then((json) => setData(json));
  }, [userId]);

  return <div>{data ? data.name : 'Loading...'}</div>;
}
