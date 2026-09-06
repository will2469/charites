import React, { useEffect, useState } from 'react';

export const UserProfileView = ({ firstName, lastName }: { firstName: string; lastName: string }) => {
  const [fullName, setFullName] = useState('');

  useEffect(() => {
    setFullName(`${firstName} ${lastName}`);
  }, [firstName, lastName]);

  return <div>Full Name: {fullName}</div>;
};
