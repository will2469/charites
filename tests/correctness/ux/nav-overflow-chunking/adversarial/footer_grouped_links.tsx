export function FooterGroupedLinks() {
  return (
    <footer className="grid grid-cols-2 gap-8">
      <div>
        <h4>Company</h4>
        <a href="/about">About</a>
        <a href="/careers">Careers</a>
        <a href="/press">Press</a>
        <a href="/news">News</a>
      </div>
      <div>
        <h4>Legal</h4>
        <a href="/terms">Terms</a>
        <a href="/privacy">Privacy</a>
        <a href="/cookies">Cookies</a>
        <a href="/licenses">Licenses</a>
      </div>
    </footer>
  );
}
