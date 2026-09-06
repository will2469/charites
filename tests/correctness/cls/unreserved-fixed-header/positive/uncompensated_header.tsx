import React from 'react';

export function UncompensatedHeaderLayout() {
  return (
    <div>
      <header className="fixed top-0 left-0 w-full h-16 bg-background z-50">
        <nav>Navbar</nav>
      </header>
      <main>
        <h1>Welcome to Dashboard</h1>
        <p>Main content without top padding gets obscured behind fixed header.</p>
      </main>
    </div>
  );
}
