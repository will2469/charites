import React from 'react';

export const UserProfileView = ({ firstName, lastName }: { firstName: string; lastName: string }) => {
  const fullName = `${firstName} ${lastName}`;

  return <div>Full Name: {fullName}</div>;
};
