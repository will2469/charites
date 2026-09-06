export function HandledInvoiceList({ invoices }: { invoices: any[] }) {
  return (
    <div className="space-y-3">
      <h2 className="text-lg font-bold">Daftar Tagihan</h2>
      {invoices.length === 0 ? (
        <EmptyState
          title="Belum Ada Tagihan"
          description="Buat tagihan pertama Anda untuk memulai."
        />
      ) : (
        <List items={invoices.map(inv => <InvoiceRow key={inv.id} data={inv} />)} />
      )}
    </div>
  );
}
