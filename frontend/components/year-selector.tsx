"use client";

import { cn } from "@/lib/utils";
import { useSeason } from "@/hooks/SeasonContext";

export function YearSelector() {
  const { selectedYear, setSelectedYear, availableYears } = useSeason();

  return (
    <div className="border-b border-border/40 bg-secondary/30">
      <div className="container mx-auto px-4 md:px-6">
        <div className="flex items-center gap-1 py-2">
          <span className="mr-3 text-sm font-medium text-muted-foreground">Season</span>
          <div className="flex gap-1 rounded-lg bg-muted/50 p-1">
            {availableYears.map((year) => (
              <button
                key={year}
                onClick={() => setSelectedYear(year)}
                className={cn(
                  "rounded-md px-4 py-1.5 text-sm font-semibold transition-all duration-200",
                  selectedYear === year
                    ? "bg-primary text-primary-foreground shadow-sm"
                    : "text-muted-foreground hover:text-foreground hover:bg-muted"
                )}
              >
                {year}
              </button>
            ))}
          </div>
        </div>
      </div>
    </div>
  );
}
