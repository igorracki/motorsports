import { AlertCircle } from "lucide-react";
import { CircuitMap } from "@/components/features/circuit-map";
import { TrackStats } from "@/components/features/track-stats";
import { LoadingSpinner } from "@/components/ui/LoadingSpinner";
import { cn } from "@/lib/utils";
import type { RaceWeekend, Circuit } from "@/types/f1";
import type { RaceStatus } from "@/lib/race-utils";
import { formatRaceRange } from "@/lib/date-utils";

interface TrackOverviewProps {
  raceWeekend: RaceWeekend;
  year: number;
  circuit: Circuit | null;
  status: RaceStatus;
  loading: boolean;
  error: string | null;
  onRetry: () => void;
}

export function TrackOverview({
  raceWeekend,
  year,
  circuit,
  status,
  loading,
  error,
  onRetry
}: TrackOverviewProps) {
  return (
    <section className="mb-8">
      <div className="grid gap-8 lg:grid-cols-2 items-stretch">
        <div className="flex flex-col">
          <div className={cn(
            "relative overflow-hidden rounded-2xl border bg-card p-6 h-full flex items-center justify-center min-h-[300px]",
            status === "ongoing" ? "border-primary/60" : "border-border/50"
          )}>
            {loading ? (
              <LoadingSpinner label="Loading track map..." />
            ) : error ? (
              <div className="flex flex-col items-center gap-2 text-muted-foreground">
                <AlertCircle className="h-8 w-8 opacity-20" />
                <p className="text-sm">{error}</p>
                <button 
                  onClick={onRetry}
                  className="text-xs text-primary hover:underline mt-2"
                >
                  Try Again
                </button>
              </div>
            ) : (
              <CircuitMap 
                layout={circuit?.layout} 
                rotation={circuit?.rotation} 
                className="w-full h-full max-w-md" 
              />
            )}
          </div>
        </div>

        <div className="flex flex-col gap-6 h-full">
          <div>
            <h2 className="mb-2 text-2xl font-bold sm:text-4xl">
              {raceWeekend.name}
            </h2>
            <div className="flex flex-wrap items-center gap-x-4 gap-y-2 text-muted-foreground">
              <p className="font-medium">
                {raceWeekend.location}, {raceWeekend.country}
              </p>
              <div className="h-4 w-px bg-border/50 hidden sm:block" />
              <p className={cn(
                "font-semibold",
                status === "ongoing" ? "text-primary" : "text-accent"
              )}>
                {formatRaceRange(raceWeekend.startDateUTCMS, raceWeekend.endDateUTCMS)}, {year}
              </p>
            </div>
          </div>
          <div className="flex-1">
            {loading ? (
              <div className="rounded-xl border border-border/50 bg-card/50 p-8 h-full flex items-center justify-center">
                <LoadingSpinner size="sm" label="Loading track stats..." />
              </div>
            ) : error ? (
              <div className="rounded-xl border border-border/50 bg-card/50 p-8 h-full flex items-center justify-center text-center">
                <p className="text-sm text-muted-foreground">{error}</p>
              </div>
            ) : (
              circuit && <TrackStats stats={circuit} />
            )}
          </div>
        </div>
      </div>
    </section>
  );
}
