import { useState } from "react";

// Positive fixture: Destructive button sets visibility for an informational dialog lacking confirmation action
export function InfoModalViolation() {
  const [isOpen, setIsOpen] = useState(false);

  return (
    <div>
      <Button variant="destructive" onClick={() => setIsOpen(true)}>
        Hapus Akun
      </Button>

      <Dialog open={isOpen}>
        <div>Pemberitahuan: Akun pengguna akan dihapus oleh sistem.</div>
      </Dialog>
    </div>
  );
}
