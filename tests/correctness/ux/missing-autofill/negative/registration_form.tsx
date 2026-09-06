import React from 'react';

export function RegistrationForm() {
  return (
    <form className="space-y-4">
      <input
        type="text"
        name="fullname"
        autoComplete="name"
        placeholder="Nama Lengkap"
      />
      <input
        type="tel"
        name="mobile"
        autoComplete="tel"
        placeholder="Nomor HP"
      />
      <input
        type="password"
        name="new_password"
        autoComplete="new-password"
        placeholder="Kata Sandi Baru"
      />
    </form>
  );
}
