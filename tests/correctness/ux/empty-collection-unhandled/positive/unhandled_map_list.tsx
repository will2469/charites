export function InvoiceList({ invoices }: { invoices: any[] }) {
  return (
    <div className="space-y-3">
      <h2 className="text-lg font-bold">Daftar Tagihan</h2>
      <List items={invoices.map(inv => <InvoiceRow key={inv.id} data={inv} />)} />
    </div>
  );
}
