import React from 'react';

export function InFlowHeaderLayout() {
  return (
    <div>
      <header className="w-full h-16 bg-background border-b">
        <nav>Navbar</nav>
      </header>
      <main>
        <h1>Welcome to Normal Flow</h1>
        <p>In-flow relative header naturally pushes main content down without shift.</p>
      </main>
    </div>
  );
}
