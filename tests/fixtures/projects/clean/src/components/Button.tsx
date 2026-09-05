import React from 'react';

export const Button: React.FC<{ label: string }> = ({ label }) => {
  return (
    <button className="bg-primary text-white px-4 py-2 rounded-md font-medium">
      {label}
    </button>
  );
};
