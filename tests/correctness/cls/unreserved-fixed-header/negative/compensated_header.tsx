import React from 'react';

export function CompensatedHeaderLayout() {
  return (
    <div>
      <header className="fixed top-0 left-0 w-full h-16 bg-background z-50">
        <nav>Navbar</nav>
      </header>
      <main className="pt-16">
        <h1>Welcome to Dashboard</h1>
        <p>Main content has matching top padding compensation.</p>
      </main>
    </div>
  );
}

export function SpacerHeaderLayout() {
  return (
    <div>
      <header className="fixed top-0 left-0 w-full h-16 bg-background z-50">
        <nav>Navbar</nav>
      </header>
      <div className="h-16" aria-hidden="true" />
      <main>
        <h1>Welcome to Dashboard</h1>
      </main>
    </div>
  );
}
