"use client";

import { useEffect, useState } from "react";
import { RaceCard } from "@/components/features/race-card";
import type { RaceWeekend } from "@/types/f1";
import { getRaceStatus } from "@/lib/race-utils";

interface RacesGridProps {
  raceWeekends: RaceWeekend[];
  year: number;
}

interface RacesGridProps {
  raceWeekends: RaceWeekend[];
  year: number;
  serverTime?: number;
}

export function RacesGrid({ raceWeekends, year, serverTime = 0 }: RacesGridProps) {
  const [now, setNow] = useState(serverTime);

  useEffect(() => {
    // Sync with local clock after mount to avoid hydration mismatch
    const handle = requestAnimationFrame(() => setNow(Date.now()));
    return () => cancelAnimationFrame(handle);
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
        const status = getRaceStatus(year, raceWeekend.round, raceWeekend, now);

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
