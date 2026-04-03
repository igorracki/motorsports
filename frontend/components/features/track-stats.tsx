"use client";

import type { Circuit } from "@/types/f1";
import { formatDayMonth } from "@/lib/date-utils";
import { cn } from "@/lib/utils";

interface TrackStatsProps {
  stats: Circuit;
}

export function TrackStats({ stats }: TrackStatsProps) {
  const items: { label: string; value: string | number | null; isDate?: boolean; isPlaceholder?: boolean }[] = [
    { label: "Length", value: stats.lengthKM ? `${stats.lengthKM.toFixed(1)} km` : "TBC" },
    { label: "Corners", value: stats.corners ?? "TBC" },
    { label: "Top Speed", value: stats.maxSpeedKmh ? `${Math.round(stats.maxSpeedKmh)} km/h` : "TBC" },
    { label: "Min Altitude", value: stats.minAltitudeM !== undefined ? `${Math.round(stats.minAltitudeM)} m` : "TBC" },
    { label: "Max Altitude", value: stats.maxAltitudeM !== undefined ? `${Math.round(stats.maxAltitudeM)} m` : "TBC" },
    { label: "Race Date", value: formatDayMonth(stats.eventDateMS), isDate: true },
  ];

  const gridItems = [...items];
  while (gridItems.length < 6) {
    gridItems.push({ label: "", value: "", isPlaceholder: true });
  }

  return (
    <div className="rounded-2xl border border-border/50 bg-card overflow-hidden shadow-sm h-full flex flex-col">
      <div className="border-b border-border/50 bg-secondary/30 px-6 py-4">
        <h3 className="text-sm font-semibold uppercase tracking-wider text-muted-foreground">
          Circuit Information
        </h3>
      </div>
      <div className="grid grid-cols-2 lg:grid-cols-3 flex-1 items-stretch">
        {gridItems.map((item, idx) => {
          const isPlaceholder = item.isPlaceholder;
          return (
            <div
              key={idx}
              className={cn(
                "flex flex-col justify-center gap-1 p-6 transition-colors hover:bg-secondary/5 border-border/30",
                // Right border logic
                (idx + 1) % 3 !== 0 && "lg:border-r",
                (idx + 1) % 2 !== 0 && "max-lg:border-r",
                // Bottom border logic: 
                // For 2 columns (max-lg), last row is idx 4,5. So border-b for idx < 4
                // For 3 columns (lg), last row is idx 3,4,5. So border-b for idx < 3
                idx < 3 && "lg:border-b",
                idx < 4 && "max-lg:border-b",
                isPlaceholder && "max-lg:hidden lg:bg-secondary/[0.02]"
              )}
            >
              {!isPlaceholder && (
                <>
                  <p className="text-[10px] font-bold uppercase tracking-[0.2em] text-muted-foreground/60 mb-1">
                    {item.label}
                  </p>
                  <p className={cn(
                    "text-xl font-bold tracking-tight sm:text-2xl",
                    item.isDate ? "text-primary" : "text-foreground"
                  )}>
                    {item.value}
                  </p>
                </>
              )}
            </div>
          );
        })}
      </div>
    </div>
  );
}

