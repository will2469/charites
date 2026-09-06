import React from 'react';

export function BloatedFieldsetForm() {
  return (
    <form className="space-y-6">
      <fieldset className="space-y-4 border p-4">
        <legend>Data Identitas Diri Lengkap</legend>
        <input name="nik" placeholder="NIK" />
        <input name="nama" placeholder="Nama Lengkap" />
        <input name="tempat_lahir" placeholder="Tempat Lahir" />
        <input name="tanggal_lahir" placeholder="Tanggal Lahir" />
        <input name="agama" placeholder="Agama" />
        <input name="pekerjaan" placeholder="Pekerjaan" />
        <input name="pendidikan" placeholder="Pendidikan" />
        <input name="gol_darah" placeholder="Golongan Darah" />
      </fieldset>
      <button type="submit">Simpan</button>
    </form>
  );
}
