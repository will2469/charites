// Negative fixture: Button satisfies both R1 (disabled) and R2 (aria-busy & dynamic text)
export function CompliantOrderButton({ isPending, onPay }: { isPending: boolean; onPay: () => void }) {
  return (
    <button
      onClick={async () => {
        await api.post("/orders", {});
      }}
      disabled={isPending}
      aria-busy={isPending}
      className="bg-primary text-white px-4 py-2"
    >
      {isPending ? "Memproses..." : "Bayar Sekarang"}
    </button>
  );
}
