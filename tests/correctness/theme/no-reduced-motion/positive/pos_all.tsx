// POS-002: Transition all on body in TSX without prefers-reduced-motion
export function BodyStyles() {
  return (
    <style>{`
      body {
        transition: all 0.2s ease-in-out;
      }
    `}</style>
  );
}
