// Positive fixture: Setter is received via props from outside the component scope
export function ChildComponent({ setConfirmOpen }: { setConfirmOpen: (v: boolean) => void }) {
  return (
    <Button variant="destructive" onClick={() => setConfirmOpen(true)}>
      Hapus Akun
    </Button>
  );
}
