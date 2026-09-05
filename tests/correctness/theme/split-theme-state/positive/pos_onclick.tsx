// POS-001: Direct localStorage mutation in button handler
export function BadToggle() {
  return (
    <button onClick={() => localStorage.setItem('theme', 'dark')}>
      Toggle Theme
    </button>
  );
}
