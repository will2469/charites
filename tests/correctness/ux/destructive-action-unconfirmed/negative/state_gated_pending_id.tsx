import { useState } from "react";

// Negative fixture: Staged confirmation with pending action identity and Boolean predicate
export function TableActionList({ items }: { items: any[] }) {
  const [pendingDeleteId, setPendingDeleteId] = useState<string | null>(null);

  return (
    <div>
      {items.map((item) => (
        <Button
          key={item.id}
          variant="destructive"
          onClick={() => setPendingDeleteId(item.id)}
        >
          Hapus Item
        </Button>
      ))}

      <ConfirmDialog
        open={Boolean(pendingDeleteId)}
        onConfirm={() => deleteItem(pendingDeleteId)}
      />
    </div>
  );
}
