import React from 'react';

export function TagContainmentBait() {
  return (
    <nav aria-label="Navigasi">
      <Breadcrumb />
      <Breadcrumb />
      <BrandedButton />
      <BrandedButton />
    </nav>
  );
}

function Breadcrumb() {
  return <span>Nav</span>;
}

function BrandedButton() {
  return <button type="button">Action</button>;
}
