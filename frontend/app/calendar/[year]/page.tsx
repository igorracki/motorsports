import { Suspense } from "react";
import { MainNav } from "@/components/ui/main-nav";
import { Footer } from "@/components/ui/Footer";
import { YearSelector } from "@/components/features/year-selector";
import { RacesGrid } from "@/components/features/races-grid";
import { getServerApi } from "@/services/server-api";
import { Calendar, Flag, Zap } from "lucide-react";
import { getScheduleStats } from "@/lib/race-utils";
import { EventCardSkeleton, Skeleton } from "@/components/ui/Skeleton";
import { ErrorBoundary } from "@/components/ui/error-boundary";

interface CalendarPageProps {
  params: Promise<{ year: string }>;
}

async function CalendarContent({ year }: { year: number }) {
  const { raceRepo } = getServerApi();
  const raceWeekends = await raceRepo.getSchedule(year);
  const now = new Date().getTime();
  const stats = getScheduleStats(raceWeekends, year, now);

  return (
    <>
      <div className="mb-8">
        <h1 className="text-3xl font-bold tracking-tight md:text-4xl">
          {year} Season
        </h1>

        <div className="mt-3 flex flex-wrap items-center gap-4 text-sm text-muted-foreground">
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
              <span>{stats.ongoing} {stats.ongoing === 1 ? 'race' : 'races'} in progress</span>
            </div>
          )}
        </div>
      </div>

      <ErrorBoundary name="Races Grid">
        <RacesGrid raceWeekends={raceWeekends} year={year} serverTime={now} />
      </ErrorBoundary>
    </>
  );
}

function CalendarFallback({ year }: { year: number }) {
  return (
    <>
      <div className="mb-8">
        <h1 className="text-3xl font-bold tracking-tight md:text-4xl">
          {year} Season
        </h1>

        <div className="mt-3 flex flex-wrap items-center gap-4 text-sm text-muted-foreground">
          <div className="flex items-center gap-2">
            <Skeleton className="h-5 w-32" />
          </div>
          <div className="flex items-center gap-2">
            <Skeleton className="h-5 w-32" />
          </div>
        </div>
      </div>

      <div className="grid gap-6 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
        {[...Array(12)].map((_, i) => (
          <EventCardSkeleton key={i} />
        ))}
      </div>
    </>
  );
}

export default async function CalendarPage({ params }: CalendarPageProps) {
  const { year: yearStr } = await params;
  const year = Number(yearStr) || 2026;

  return (
    <div className="min-h-screen bg-background text-foreground flex flex-col">
      <MainNav />
      <YearSelector currentYear={year} />

      <main className="container mx-auto px-4 py-8 md:px-6 flex-1">
        <Suspense fallback={<CalendarFallback year={year} />}>
          <CalendarContent year={year} />
        </Suspense>
      </main>

      <Footer />
    </div>
  );
}
