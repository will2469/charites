// Positive fixture: Async handler turns on loading flag without finally/reset on error
export function LoadButton() {
  return (
    <button
      onClick={async () => {
        setLoading(true);
        try {
          await api.fetchUsers();
        } catch (e) {
          console.error(e);
        }
      }}
      className="bg-primary text-white p-2"
    >
      Muat Data
    </button>
  );
}
