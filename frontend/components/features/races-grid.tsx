"use client";

import { RaceCard } from "@/components/features/race-card";
import type { RaceWeekend } from "@/types/f1";

interface RacesGridProps {
  raceWeekends: RaceWeekend[];
  year: number;
}

import { getRaceStatus } from "@/lib/race-utils";

export function RacesGrid({ raceWeekends, year }: RacesGridProps) {
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
      {raceWeekends.map((raceWeekend) => (
        <RaceCard
          key={raceWeekend.round}
          raceWeekend={raceWeekend}
          year={year}
          status={getRaceStatus(year, raceWeekend.round, raceWeekend)}
        />
      ))}
    </div>
  );
}
