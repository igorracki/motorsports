"use client";

import { useEffect, useState } from "react";
import { useSeason } from "@/hooks/SeasonContext";
import { MainNav } from "@/components/main-nav";
import { YearSelector } from "@/components/year-selector";
import { RacesGrid } from "@/components/races-grid";
import { f1Api } from "@/services/f1-api";
import { RaceWeekend } from "@/types/f1";
import { Calendar, Flag, Zap } from "lucide-react";
import { Skeleton, EventCardSkeleton } from "@/components/ui/Skeleton";

import { getScheduleStats } from "@/lib/race-utils";

export function CalendarView() {
  const { selectedYear } = useSeason();
  const [raceWeekends, setRaceWeekends] = useState<RaceWeekend[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    async function loadSchedule() {
      setLoading(true);
      try {
        const schedule = await f1Api.getSchedule(selectedYear);
        setRaceWeekends(schedule);
      } catch (error) {
        console.error("Failed to load schedule:", error);
      } finally {
        setLoading(false);
      }
    }
    loadSchedule();
  }, [selectedYear]);

  // Derived stats
  const stats = getScheduleStats(raceWeekends, selectedYear);

  return (
    <div className="min-h-screen bg-background text-foreground">
      <MainNav />
      <YearSelector />

      <main className="container mx-auto px-4 py-8 md:px-6">
        <div className="mb-8">
          <h1 className="text-3xl font-bold tracking-tight md:text-4xl">
            {selectedYear} Season
          </h1>
          
          <div className="mt-3 flex flex-wrap items-center gap-4 text-sm text-muted-foreground">
            {loading ? (
              <>
                <Skeleton className="h-5 w-32" />
                <Skeleton className="h-5 w-32" />
              </>
            ) : (
              <>
                <div className="flex items-center gap-2">
                  <Calendar className="h-4 w-4 text-primary" />
                  <span>{stats.total} races scheduled</span>
                </div>
                {stats.completed > 0 && (
                  <div className="flex items-center gap-2">
                    <Flag className="h-4 w-4 text-accent" />
                    <span>{stats.completed} races completed</span>
                  </div>
                )}
                {stats.ongoing > 0 && (
                  <div className="flex items-center gap-2">
                    <Zap className="h-4 w-4 text-primary" />
                    <span>{stats.ongoing} race in progress</span>
                  </div>
                )}
              </>
            )}
          </div>
        </div>

        {loading ? (
          <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
            {[1, 2, 3, 4, 5, 6, 7, 8].map((i) => (
              <EventCardSkeleton key={i} />
            ))}
          </div>
        ) : (
          <RacesGrid raceWeekends={raceWeekends} year={selectedYear} />
        )}

        <footer className="mt-16 border-t border-border/40 pt-8 text-center">
          <p className="text-sm text-muted-foreground">
            F1 Data Hub {selectedYear} - Official F1 Timing and Results
          </p>
        </footer>
      </main>
    </div>
  );
}
