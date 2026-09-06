export function CheckoutCart({ cartTotal }: { cartTotal: number }) {
  return (
    <div className="mt-4">
      <button disabled={cartTotal < 50000} className="bg-primary text-white px-4 py-2 rounded">
        Checkout
      </button>
    </div>
  );
}
