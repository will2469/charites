// NEG-002: Mitigated with prefers-reduced-motion: reduce
export function AccessibleStyles() {
  return (
    <style>{`
      body {
        transition: background-color 0.3s ease;
      }
      @media (prefers-reduced-motion: reduce) {
        body {
          transition: none;
        }
      }
    `}</style>
  );
}
