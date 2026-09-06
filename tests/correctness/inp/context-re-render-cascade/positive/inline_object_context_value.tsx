import React, { createContext } from 'react';

export const AuthContext = createContext<any>(null);

export function AuthProvider({ user, isAuthenticated, children }: any) {
  return (
    <AuthContext.Provider value={{ user, isAuthenticated }}>
      {children}
    </AuthContext.Provider>
  );
}
