"use client";

import Link from "next/link";
import Image from "next/image";
import { cn } from "@/lib/utils";
import { StatusBadge } from "./ui/StatusBadge";
import { RaceWeekend } from "@/types/f1";

interface RaceCardProps {
  raceWeekend: RaceWeekend;
  year: number;
  status: "completed" | "ongoing" | "upcoming";
}

export function RaceCard({
  raceWeekend,
  year,
  status,
}: RaceCardProps) {
  const isOngoing = status === "ongoing";
  const isUpcoming = status === "upcoming";

  // Use round as the identifier in the URL
  return (
    <Link href={`/race-weekend/${year}/${raceWeekend.round}`}>
      <div
        className={cn(
          "group relative cursor-pointer overflow-hidden rounded-2xl border bg-card p-5 transition-all duration-300 hover:shadow-lg hover:shadow-primary/5",
          isOngoing
            ? "border-primary/60 shadow-md shadow-primary/10"
            : isUpcoming
              ? "border-border/50 opacity-60 hover:opacity-80"
              : "border-border/50 hover:border-primary/30"
        )}
      >
        <StatusBadge status={status} className="absolute right-3 top-3" />

        <div className="mb-4 flex items-start justify-between">
          <div className="flex-1 pr-16">
            <h3
              className={cn(
                "text-lg font-bold tracking-tight transition-colors",
                isOngoing
                  ? "text-primary"
                  : "text-foreground group-hover:text-primary"
              )}
            >
              {raceWeekend.name}
            </h3>
          </div>
          <div className="absolute right-5 top-12">
            {raceWeekend.countryCode && (
              <Image
                src={`https://flagcdn.com/w80/${raceWeekend.countryCode.toLowerCase()}.png`}
                alt={`${raceWeekend.country} flag`}
                width={36}
                height={24}
                className="rounded-sm object-cover shadow-sm"
                unoptimized
              />
            )}
          </div>
        </div>

        <div className="mb-4 space-y-1.5">
          <p className="text-sm text-muted-foreground">{raceWeekend.location}</p>
          <p
            className={cn(
              "text-sm font-medium",
              isOngoing ? "text-primary" : "text-accent"
            )}
          >
            {raceWeekend.startDate || "TBC"}
          </p>
        </div>

        <div className="flex flex-wrap gap-1.5">
          {raceWeekend.sessions.map((session, index) => (
            <span
              key={index}
              className={cn(
                "inline-flex h-7 w-8 items-center justify-center rounded-md text-xs font-bold transition-colors",
                isOngoing
                  ? "bg-primary/20 text-primary"
                  : "bg-secondary text-secondary-foreground group-hover:bg-primary/20 group-hover:text-primary"
              )}
            >
              {session.type}
            </span>
          ))}
        </div>
      </div>
    </Link>
  );
}
