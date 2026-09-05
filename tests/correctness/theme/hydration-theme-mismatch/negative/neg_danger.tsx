// NEG-002: TSX with dangerouslySetInnerHTML theme script
export function SafeLayout({ children }: { children: React.ReactNode }) {
  return (
    <html>
      <head>
        <meta charSet="utf-8" />
        <script
          dangerouslySetInnerHTML={{
            __html: "document.documentElement.classList.add(localStorage.getItem('theme') || 'dark');",
          }}
        />
      </head>
      <body>{children}</body>
    </html>
  );
}
