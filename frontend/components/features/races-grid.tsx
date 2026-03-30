"use client";

import { useEffect, useState } from "react";
import { RaceCard } from "@/components/features/race-card";
import type { RaceWeekend } from "@/types/f1";
import { getRaceStatus } from "@/lib/race-utils";

interface RacesGridProps {
  raceWeekends: RaceWeekend[];
  year: number;
}

export function RacesGrid({ raceWeekends, year }: RacesGridProps) {
  const [isMounted, setIsMounted] = useState(false);

  useEffect(() => {
    // eslint-disable-next-line react-hooks/set-state-in-effect
    setIsMounted(true);
  }, []);

  if (raceWeekends.length === 0) {
    return (
      <div className="flex min-h-[300px] items-center justify-center">
        <p className="text-muted-foreground">
          No races scheduled for this season.
        </p>
      </div>
    );
  }

  return (
    <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
      {raceWeekends.map((raceWeekend) => {
        // Hydration guard: Force all races to "upcoming" status on the server
        // and initial client render to match the server perfectly.
        const status = isMounted
          ? getRaceStatus(year, raceWeekend.round, raceWeekend)
          : "upcoming";

        return (
          <RaceCard
            key={raceWeekend.round}
            raceWeekend={raceWeekend}
            year={year}
            status={status}
          />
        );
      })}
    </div>
  );
}
