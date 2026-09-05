// POS-002: Missing color-scheme in data-theme="dark"
export function ThemeStyles() {
  return (
    <style>{`
      [data-theme="dark"] {
        --bg: #121212;
      }
    `}</style>
  );
}
