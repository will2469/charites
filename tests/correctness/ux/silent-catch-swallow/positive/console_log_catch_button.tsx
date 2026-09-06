// Positive fixture: Catch block logging to console without user feedback or rethrow
export function ProfileForm({ data }: { data: any }) {
  return (
    <button
      onClick={async () => {
        try {
          await api.updateProfile(data);
        } catch (e) {
          console.error(e);
        }
      }}
      className="btn"
    >
      Simpan Profil
    </button>
  );
}
