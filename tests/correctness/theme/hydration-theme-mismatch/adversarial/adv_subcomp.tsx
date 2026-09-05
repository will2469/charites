// ADV-001: Normal component without head tag (no false positive)
export function UserProfile() {
  return (
    <div className="p-6 bg-card text-foreground rounded-lg">
      <h2>User Profile</h2>
      <p>Content goes here.</p>
    </div>
  );
}
