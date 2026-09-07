import { useState } from "react";

// Positive fixture: Handler sets an unrelated state but directly executes deleteUser
export function MixedUnrelatedSetter({ id }: { id: string }) {
  const [selectedTab, setSelectedTab] = useState("all");

  return (
    <Button
      variant="destructive"
      onClick={() => {
        setSelectedTab("all");
        deleteUser(id);
      }}
    >
      Hapus Pengguna
    </Button>
  );
}
