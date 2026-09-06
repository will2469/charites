// Positive fixture: Async submit button without reentry guard or feedback
export function OrderButton({ data }: { data: any }) {
  return (
    <button
      onClick={async () => {
        await api.post("/orders", data);
      }}
      className="bg-primary text-white px-4 py-2"
    >
      Bayar Sekarang
    </button>
  );
}
