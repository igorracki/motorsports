"use client";

import type { Circuit } from "@/types/f1";

interface TrackStatsProps {
  stats: Circuit;
}

export function TrackStats({ stats }: TrackStatsProps) {
  return (
    <div className="rounded-2xl border border-border/50 bg-card p-5 sm:p-6">
      <h3 className="mb-4 text-sm font-semibold uppercase tracking-wider text-muted-foreground">
        Circuit Information
      </h3>
      <div className="flex flex-wrap items-start gap-x-8 gap-y-4 sm:gap-x-12">
        <div className="min-w-[80px]">
          <p className="text-2xl font-bold text-foreground">{stats.lengthKM} km</p>
          <p className="text-xs text-muted-foreground">Length</p>
        </div>
        <div className="min-w-[60px]">
          <p className="text-2xl font-bold text-foreground">{stats.corners}</p>
          <p className="text-xs text-muted-foreground">Corners</p>
        </div>
        <div className="min-w-[80px]">
          <p className="text-2xl font-bold text-primary">{stats.eventDate || "TBC"}</p>
          <p className="text-xs text-muted-foreground">Race Date</p>
        </div>
      </div>
    </div>
  );
}
