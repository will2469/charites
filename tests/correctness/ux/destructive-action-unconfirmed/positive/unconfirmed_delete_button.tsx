// Positive fixture: Delete button triggering account removal directly without confirmation
export function UserRow({ user }: { user: any }) {
  return (
    <div className="flex justify-between items-center p-2">
      <span>{user.name}</span>
      <button
        onClick={() => deleteUser(user.id)}
        className="bg-destructive text-destructive-foreground px-3 py-1 rounded"
      >
        Hapus Akun
      </button>
    </div>
  );
}
