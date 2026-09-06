// Positive fixture: Form handler sets email error, but email error is never rendered in JSX
export function LoginForm() {
  return (
    <form onSubmit={(e) => { e.preventDefault(); setEmailError("Format email tidak valid"); }} className="space-y-4">
      <input type="email" placeholder="Email Anda" className="border px-3 py-2" />
      <button type="submit" className="bg-primary text-white">Masuk</button>
    </form>
  );
}
