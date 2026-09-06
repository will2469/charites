export function CompliantCheckoutCart({ cartTotal }: { cartTotal: number }) {
  return (
    <div className="mt-4">
      <button
        disabled={cartTotal < 50000}
        aria-describedby="min-order-hint"
        className="bg-primary text-white px-4 py-2 rounded"
      >
        Checkout
      </button>
      <p id="min-order-hint" className="text-xs text-muted-foreground mt-1">
        Minimum belanja Rp 50.000 untuk melanjutkan.
      </p>
    </div>
  );
}
