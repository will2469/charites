// Negative fixture: Catch block displays user feedback via toast.error
export function CompliantProfileForm({ data }: { data: any }) {
  return (
    <button
      onClick={async () => {
        try {
          await api.updateProfile(data);
        } catch (e) {
          toast.error("Gagal memperbarui profil. Coba lagi.");
        }
      }}
      className="btn"
    >
      Simpan Profil
    </button>
  );
}
