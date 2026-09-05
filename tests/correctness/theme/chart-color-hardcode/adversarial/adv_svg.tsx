// ADV-001: Standard SVG path and non-chart element with fill
export function IconWrapper() {
  return (
    <div>
      <svg viewBox="0 0 24 24">
        <path d="M0 0h24v24H0z" fill="none" />
      </svg>
      <div className="bg-card text-foreground">Content</div>
    </div>
  );
}
