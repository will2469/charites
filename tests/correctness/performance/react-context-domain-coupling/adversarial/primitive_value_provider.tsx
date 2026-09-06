import React, { createContext } from 'react';

export const ThemeContext = createContext<string>('light');

export const ThemeProvider = ({ children }: { children: React.ReactNode }) => {
  return (
    <ThemeContext.Provider value="dark">
      {children}
    </ThemeContext.Provider>
  );
};
