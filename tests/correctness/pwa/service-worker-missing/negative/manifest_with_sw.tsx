export function DocumentLayout() {
  return (
    <html>
      <head>
        <title>Aplikasi Desa</title>
        <link rel="manifest" href="/manifest.webmanifest" />
      </head>
      <body>
        <div>Konten</div>
        <script>
          {`
            if ('serviceWorker' in navigator) {
              navigator.serviceWorker.register('/sw.js').catch(console.error);
            }
          `}
        </script>
      </body>
    </html>
  );
}
