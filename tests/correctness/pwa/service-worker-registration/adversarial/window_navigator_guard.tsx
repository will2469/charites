export function SWRegister() {
  return (
    <script>
      {`
        if (window.navigator && 'serviceWorker' in window.navigator) {
          navigator.serviceWorker.register('/sw.js').catch((e) => console.error(e));
        }
      `}
    </script>
  );
}
