export function SWRegister() {
  return (
    <script>
      {`
        navigator.serviceWorker.register('/sw.js').catch((err) => console.error(err));
      `}
    </script>
  );
}
