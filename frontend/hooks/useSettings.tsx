"use client";

import React, { createContext, useContext, useEffect, useState } from "react";

interface SettingsContextType {
  useBrowserTime: boolean;
  setUseBrowserTime: (value: boolean) => void;
}

const SettingsContext = createContext<SettingsContextType | undefined>(undefined);

export const SettingsProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const [useBrowserTime, setUseBrowserTime] = useState<boolean>(true);
  const [isLoaded, setIsLoaded] = useState(false);

  useEffect(() => {
    const saved = localStorage.getItem("useBrowserTime");
    // Use timeout to prevent synchronous state updates during effect phase
    const handle = setTimeout(() => {
      if (saved !== null) {
        setUseBrowserTime(JSON.parse(saved));
      }
      setIsLoaded(true);
    }, 0);
    return () => clearTimeout(handle);
  }, []);

  useEffect(() => {
    if (isLoaded) {
      localStorage.setItem("useBrowserTime", JSON.stringify(useBrowserTime));
    }
  }, [useBrowserTime, isLoaded]);

  return (
    <SettingsContext.Provider value={{ useBrowserTime, setUseBrowserTime }}>
      {children}
    </SettingsContext.Provider>
  );
};

export const useSettings = () => {
  const context = useContext(SettingsContext);
  if (context === undefined) {
    throw new Error("useSettings must be used within a SettingsProvider");
  }
  return context;
};
