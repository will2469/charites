// NEG-002: Pure media-query based dark mode
export function MediaOnly() {
  return (
    <style>{`
      :root {
        --bg-main: #ffffff;
      }
      @media (prefers-color-scheme: dark) {
        :root {
          color-scheme: dark;
          --bg-main: #09090b;
        }
      }
    `}</style>
  );
}
