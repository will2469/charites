import React from 'react';

export function ShortLoginForm() {
  return (
    <form className="space-y-4 max-w-sm">
      <input name="username" placeholder="Username" />
      <input type="password" name="password" placeholder="Password" />
      <button type="submit">Masuk</button>
    </form>
  );
}
