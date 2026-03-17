"use client";

import React, { createContext, useContext, useMemo, useCallback, useState, ReactNode, useEffect, Suspense, useTransition } from "react";
import { useSearchParams, useRouter, usePathname } from "next/navigation";

interface SeasonContextType {
  selectedYear: number;
  setSelectedYear: (year: number) => void;
  availableYears: number[];
  isPending: boolean;
}

const SeasonContext = createContext<SeasonContextType | undefined>(undefined);

function URLSync({ onSync }: { onSync: (year: number) => void }) {
  const searchParams = useSearchParams();
  
  useEffect(() => {
    const yearParam = searchParams.get("year");
    if (yearParam) {
      const year = parseInt(yearParam, 10);
      if (!isNaN(year)) {
        onSync(year);
      }
    } else {
      // Default to 2026 if no param, but only if we aren't already there
      onSync(2026);
    }
  }, [searchParams, onSync]);

  return null;
}

export function SeasonProvider({ children }: { children: ReactNode }) {
  const router = useRouter();
  const pathname = usePathname();
  const [selectedYear, setSelectedYearState] = useState(2026);
  const [isPending, startTransition] = useTransition();
  
  const availableYears = useMemo(() => [2026, 2025], []);

  const setSelectedYear = useCallback((year: number) => {
    // 1. Update local state immediately for snappy UI
    setSelectedYearState(year);
    
    // 2. Update URL as a transition to avoid blocking UI or triggering top-level suspense
    startTransition(() => {
      const params = new URLSearchParams(window.location.search);
      params.set("year", year.toString());
      router.replace(`${pathname}?${params.toString()}`, { scroll: false });
    });
  }, [pathname, router]);

  const value = useMemo(() => ({
    selectedYear,
    setSelectedYear,
    availableYears,
    isPending
  }), [selectedYear, setSelectedYear, availableYears, isPending]);

  return (
    <SeasonContext.Provider value={value}>
      {/* 
          Suspense boundary here captures any suspension from useSearchParams 
          without unmounting the 'children' below.
      */}
      <Suspense fallback={null}>
        <URLSync onSync={setSelectedYearState} />
      </Suspense>
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
