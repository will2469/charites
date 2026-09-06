export function FixedBottomInForm() {
  return (
    <form className="h-screen flex flex-col justify-between">
      <div className="p-4 space-y-4">
        <input type="text" placeholder="Nama Lengkap" />
        <input type="email" placeholder="Alamat Surel" />
        <textarea placeholder="Pesan Anda" />
      </div>
      <div className="fixed bottom-0 inset-x-0 p-4 bg-surface border-t">
        <button type="submit" className="w-full bg-primary text-white py-3 rounded">
          Kirim Berkas
        </button>
      </div>
    </form>
  );
}
