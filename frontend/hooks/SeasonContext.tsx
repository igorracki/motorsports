"use client";

import React, { createContext, useContext, useState, useEffect, ReactNode } from "react";
import { useSearchParams, useRouter, usePathname } from "next/navigation";

interface SeasonContextType {
  selectedYear: number;
  setSelectedYear: (year: number) => void;
  availableYears: number[];
}

const SeasonContext = createContext<SeasonContextType | undefined>(undefined);

export function SeasonProvider({ children }: { children: ReactNode }) {
  const router = useRouter();
  const pathname = usePathname();
  const searchParams = useSearchParams();
  
  // Available years for the F1 Calendar
  const availableYears = [2026, 2025];
  
  // Initialize from URL search params if present, otherwise default to latest
  const [selectedYear, setSelectedYearState] = useState<number>(() => {
    const yearParam = searchParams.get("year");
    const year = yearParam ? parseInt(yearParam, 10) : 2026;
    return availableYears.includes(year) ? year : 2026;
  });

  // Keep state in sync with URL
  const setSelectedYear = (year: number) => {
    setSelectedYearState(year);
    const params = new URLSearchParams(searchParams.toString());
    params.set("year", year.toString());
    router.push(`${pathname}?${params.toString()}`);
  };

  return (
    <SeasonContext.Provider value={{ selectedYear, setSelectedYear, availableYears }}>
      {children}
    </SeasonContext.Provider>
  );
}

export function useSeason() {
  const context = useContext(SeasonContext);
  if (context === undefined) {
    throw new Error("useSeason must be used within a SeasonProvider");
  }
  return context;
}
