import { useState } from "react";

// Negative fixture: Staged confirmation with custom setter naming resolved via useState AST binding
export function CustomSetterPattern({ item }: { item: any }) {
  const [pendingAction, updatePendingAction] = useState<string | null>(null);

  return (
    <div>
      <Button
        variant="destructive"
        onClick={() => updatePendingAction(item.id)}
      >
        Hapus Akun
      </Button>

      <ActionApprovalDialog
        open={pendingAction !== null}
        onConfirm={() => executeAction(pendingAction)}
      />
    </div>
  );
}
