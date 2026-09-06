export function CompactNavbar() {
  return (
    <nav className="flex items-center gap-6">
      <a href="/dashboard">Dashboard</a>
      <a href="/analytics">Analytics</a>
      <a href="/reports">Reports</a>
      <a href="/settings">Settings</a>
    </nav>
  );
}
