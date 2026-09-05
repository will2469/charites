// NEG-002: Valid global color-scheme: light dark
export function ValidRoot() {
  return (
    <style>{`
      :root {
        color-scheme: light dark;
      }
      .dark {
        --bg: #09090b;
      }
    `}</style>
  );
}
