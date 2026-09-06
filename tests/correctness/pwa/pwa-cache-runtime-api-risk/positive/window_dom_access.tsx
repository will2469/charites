export function WorkerScript() {
  return (
    <script>
      {`
        self.addEventListener('install', (event) => {
          window.location.reload();
          document.getElementById('root');
        });
      `}
    </script>
  );
}
