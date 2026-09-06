// Negative fixture: Form handler sets email error, rendered with role="alert"
export function CompliantLoginForm() {
  return (
    <form onSubmit={(e) => { e.preventDefault(); setEmailError("Format email tidak valid"); }} className="space-y-4">
      <input type="email" placeholder="Email Anda" className="border px-3 py-2" />
      <p role="alert" className="text-sm text-destructive font-medium">Format email tidak valid</p>
      <button type="submit" className="bg-primary text-white">Masuk</button>
    </form>
  );
}
