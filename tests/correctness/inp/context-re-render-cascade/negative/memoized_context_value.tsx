import React, { createContext, useMemo } from 'react';

export const AuthContext = createContext<any>(null);

export function AuthProvider({ user, isAuthenticated, children }: any) {
  const authValue = useMemo(() => ({ user, isAuthenticated }), [user, isAuthenticated]);

  return (
    <AuthContext.Provider value={authValue}>
      {children}
    </AuthContext.Provider>
  );
}
