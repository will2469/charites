import { useState } from "react";

// Negative fixture: Staged deletion flow with useState and ActionApprovalDialog (Issue #11 specimen)
export function FollowerList({ users, onHapus }: Props) {
  const [confirmIndex, setConfirmIndex] = useState<number | null>(null);

  return (
    <div>
      {users.map((user, index) => (
        <Button
          key={user.id}
          type="button"
          variant="destructive"
          onClick={() => setConfirmIndex(index)}
        >
          Hapus
        </Button>
      ))}

      <ActionApprovalDialog
        open={confirmIndex !== null}
        title="Hapus Pengikut"
        onConfirm={() => {
          if (confirmIndex !== null) onHapus(confirmIndex);
          setConfirmIndex(null);
        }}
        onCancel={() => setConfirmIndex(null)}
      />
    </div>
  );
}
