import React, { createContext } from 'react';

export const AppContext = createContext<any>(null);

export const AppProvider = ({ children }: { children: React.ReactNode }) => {
  const user = {};
  const cart = {};
  const theme = 'dark';
  const notifications = [];
  const activeModal = null;
  const isSidebarOpen = false;
  const locale = 'en';

  return (
    <AppContext.Provider value={{ user, cart, theme, notifications, activeModal, isSidebarOpen, locale }}>
      {children}
    </AppContext.Provider>
  );
};
