export function SmallAvatar() {
  return (
    <div className="flex items-center gap-2">
      <img src="/avatar.png" width={48} height={48} alt="User Avatar" />
      <span>User Name</span>
    </div>
  );
}
