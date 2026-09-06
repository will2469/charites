export function WorkerScript() {
  return (
    <script>
      {`
        self.addEventListener('install', (event) => {
          const windowSize = 100;
          const myDocument = "doc";
          const storageKey = "session";
        });
      `}
    </script>
  );
}
