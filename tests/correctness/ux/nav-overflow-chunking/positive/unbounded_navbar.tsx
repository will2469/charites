export function UnboundedNavbar() {
  return (
    <div role="navigation" className="flex items-center gap-6">
      <a href="/dashboard">Dashboard</a>
      <a href="/analytics">Analytics</a>
      <a href="/reports">Reports</a>
      <a href="/users">Users</a>
      <a href="/roles">Roles</a>
      <a href="/permissions">Permissions</a>
      <a href="/audit-logs">Audit Logs</a>
      <a href="/integrations">Integrations</a>
      <a href="/settings">Settings</a>
    </div>
  );
}
