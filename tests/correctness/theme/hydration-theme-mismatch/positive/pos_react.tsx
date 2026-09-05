// POS-002: TSX document head in dark theme layout without theme initialization script
export function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <html className="dark">
      <head>
        <meta charSet="utf-8" />
        <title>Dashboard</title>
      </head>
      <body>{children}</body>
    </html>
  );
}
