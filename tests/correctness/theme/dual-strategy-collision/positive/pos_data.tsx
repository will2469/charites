// POS-002: Mixing prefers-color-scheme with [data-theme="dark"]
export function ThemeBox() {
  return (
    <style>{`
      @media (prefers-color-scheme: dark) {
        :root { --card: #18181b; }
      }
      [data-theme="dark"] {
        --card: #09090b;
      }
    `}</style>
  );
}
