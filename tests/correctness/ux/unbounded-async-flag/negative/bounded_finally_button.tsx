// Negative fixture: Loading flag guaranteed reset in finally block
export function CompliantLoadButton() {
  return (
    <button
      onClick={async () => {
        setLoading(true);
        try {
          await api.fetchUsers();
        } finally {
          setLoading(false);
        }
      }}
      className="bg-primary text-white p-2"
    >
      Muat Data
    </button>
  );
}
