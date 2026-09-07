export function QRFallbackLogin({ cerberus, onCancel }: any) {
  return (
    <ChallengeSlider
      onSolve={cerberus.solveChallenge}
      onCancel={onCancel}
      error={cerberus.error}
      solved={cerberus.challengeSolved}
    />
  );
}
