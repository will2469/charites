// Negative fixture: Destructive button gated inside AlertDialogTrigger modal wrapper
export function SafeUserRow({ user }: { user: any }) {
  return (
    <div className="flex justify-between items-center p-2">
      <span>{user.name}</span>
      <AlertDialogTrigger asChild>
        <button className="bg-destructive text-destructive-foreground px-3 py-1 rounded">
          Hapus Akun
        </button>
      </AlertDialogTrigger>
    </div>
  );
}
