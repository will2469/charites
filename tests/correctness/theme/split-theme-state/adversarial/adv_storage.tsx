// ADV-001: localStorage for non-theme data (auth-token, cart)
export function AuthButton() {
  return (
    <button onClick={() => localStorage.getItem("auth-token")}>
      Check Login
    </button>
  );
}
