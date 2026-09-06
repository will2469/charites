import React, { createContext } from 'react';

export const ThemeContext = createContext<string>('light');

export function ThemeProvider({ children }: any) {
  return (
    <ThemeContext.Provider value="dark">
      {children}
    </ThemeContext.Provider>
  );
}
