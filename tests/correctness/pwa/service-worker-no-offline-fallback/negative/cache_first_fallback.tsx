export function SWWorker() {
  return (
    <script>
      {`
        self.addEventListener('fetch', (event) => {
          event.respondWith(
            caches.match(event.request).then((res) => res || fetch(event.request))
          );
        });
      `}
    </script>
  );
}
