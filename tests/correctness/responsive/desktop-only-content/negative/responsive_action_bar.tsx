export function ResponsiveActionBar() {
  return (
    <div className="flex flex-col md:flex-row items-center justify-between p-4">
      <span className="font-bold">Total: Rp 150.000</span>
      <button
        type="submit"
        className="flex w-full md:w-auto items-center justify-center px-6 py-3 bg-primary text-white rounded-lg"
      >
        Bayar Sekarang
      </button>
    </div>
  );
}
