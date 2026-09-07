import React from 'react';

export function CommentBetweenBR() {
  return (
    <div className="card-body">
      <h4>Pemberitahuan Sistem</h4>
      <br />
      {/* spacing break */}
      <br />
      <p>Harap perbarui kata sandi Anda sebelum tanggal jatuh tempo.</p>
    </div>
  );
}
