import React from 'react';

export function MultiGroupBinaryForm() {
  return (
    <form className="space-y-4">
      {/* 3 distinct binary radio pairs = 6 total radio elements, max 2 per group */}
      <fieldset>
        <legend>Status Akun</legend>
        <label><input type="radio" name="account_status" value="active" /> Aktif</label>
        <label><input type="radio" name="account_status" value="inactive" /> Non-aktif</label>
      </fieldset>

      <fieldset>
        <legend>Jenis Kelamin</legend>
        <label><input type="radio" name="gender" value="m" /> Laki-laki</label>
        <label><input type="radio" name="gender" value="f" /> Perempuan</label>
      </fieldset>

      <fieldset>
        <legend>Tipe Notifikasi</legend>
        <label><input type="radio" name="notify_type" value="sms" /> SMS</label>
        <label><input type="radio" name="notify_type" value="email" /> Email</label>
      </fieldset>
    </form>
  );
}
