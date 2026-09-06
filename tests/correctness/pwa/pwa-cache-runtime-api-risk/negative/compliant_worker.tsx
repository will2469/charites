export function WorkerScript() {
  return (
    <script>
      {`
        self.addEventListener('install', (event) => {
          event.waitUntil(
            caches.open('v1').then((cache) => cache.addAll(['/offline.html']))
          );
        });
      `}
    </script>
  );
}
