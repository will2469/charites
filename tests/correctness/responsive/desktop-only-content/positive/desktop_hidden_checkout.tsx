export function DesktopHiddenCheckout() {
  return (
    <div className="flex items-center justify-between p-4">
      <span className="font-bold">Total: Rp 150.000</span>
      <button
        type="submit"
        className="hidden md:flex items-center px-6 py-3 bg-primary text-white rounded-lg"
      >
        Bayar Sekarang
      </button>
    </div>
  );
}
