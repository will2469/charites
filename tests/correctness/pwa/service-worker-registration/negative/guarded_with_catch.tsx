export function SWRegister() {
  return (
    <script>
      {`
        if ('serviceWorker' in navigator) {
          navigator.serviceWorker.register('/sw.js')
            .then((reg) => console.log('registered', reg))
            .catch((err) => console.error(err));
        }
      `}
    </script>
  );
}
