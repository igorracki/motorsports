"use client";

import type { Circuit } from "@/types/f1";

interface TrackStatsProps {
  stats: Circuit;
}

export function TrackStats({ stats }: TrackStatsProps) {
  const elevationChange = (stats.maxAltitudeM !== undefined && stats.minAltitudeM !== undefined) 
    ? Math.round(stats.maxAltitudeM - stats.minAltitudeM) 
    : null;

  return (
    <div className="rounded-2xl border border-border/50 bg-card p-5 sm:p-6">
      <h3 className="mb-4 text-sm font-semibold uppercase tracking-wider text-muted-foreground">
        Circuit Information
      </h3>
      <div className="flex flex-wrap items-start gap-x-8 gap-y-4 sm:gap-x-12">
        <div className="min-w-[80px]">
          <p className="text-2xl font-bold text-foreground">{stats.lengthKM?.toFixed(3)} km</p>
          <p className="text-xs text-muted-foreground">Length</p>
        </div>
        <div className="min-w-[60px]">
          <p className="text-2xl font-bold text-foreground">{stats.corners}</p>
          <p className="text-xs text-muted-foreground">Corners</p>
        </div>
        {stats.maxSpeedKmh && (
          <div className="min-w-[80px]">
            <p className="text-2xl font-bold text-foreground">{Math.round(stats.maxSpeedKmh)} km/h</p>
            <p className="text-xs text-muted-foreground">Top Speed</p>
          </div>
        )}
        {elevationChange !== null && (
          <div className="min-w-[80px]">
            <p className="text-2xl font-bold text-foreground">{elevationChange} m</p>
            <p className="text-xs text-muted-foreground">Elevation</p>
          </div>
        )}
        <div className="min-w-[80px]">
          <p className="text-2xl font-bold text-primary">{stats.eventDate || "TBC"}</p>
          <p className="text-xs text-muted-foreground">Race Date</p>
        </div>
      </div>
    </div>
  );
}
