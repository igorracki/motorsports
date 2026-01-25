"use client";

import Link from "next/link";
import { cn } from "@/lib/utils";

interface EventCardProps {
  id: string;
  title: string;
  country: string;
  countryCode: string;
  location: string;
  dateFrom: string;
  dateTo: string;
  sessions: string[];
  status: "completed" | "ongoing" | "upcoming";
  year: number;
}

export function EventCard({
  id,
  title,
  country,
  countryCode,
  location,
  dateFrom,
  dateTo,
  sessions,
  status,
  year,
}: EventCardProps) {
  const isOngoing = status === "ongoing";
  const isCompleted = status === "completed";
  const isUpcoming = status === "upcoming";

  return (
    <Link href={`/event/${id}?year=${year}`}>
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
        {isOngoing && (
          <div className="absolute right-3 top-3 flex items-center gap-1.5 rounded-md bg-primary/20 px-2 py-0.5 text-xs font-semibold text-primary">
            <span className="relative flex h-2 w-2">
              <span className="absolute inline-flex h-full w-full animate-ping rounded-full bg-primary opacity-75" />
              <span className="relative inline-flex h-2 w-2 rounded-full bg-primary" />
            </span>
            LIVE
          </div>
        )}

        {isUpcoming && (
          <div className="absolute right-3 top-3 rounded-md bg-muted px-2 py-0.5 text-xs font-medium text-muted-foreground">
            UPCOMING
          </div>
        )}

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
              {title}
            </h3>
          </div>
          <div className="absolute right-5 top-12">
            <img
              src={`https://flagcdn.com/w80/${countryCode.toLowerCase()}.png`}
              alt={`${country} flag`}
              className="h-6 w-9 rounded-sm object-cover shadow-sm"
            />
          </div>
        </div>

        <div className="mb-4 space-y-1.5">
          <p className="text-sm text-muted-foreground">{location}</p>
          <p
            className={cn(
              "text-sm font-medium",
              isOngoing ? "text-primary" : "text-accent"
            )}
          >
            {dateFrom} - {dateTo}
          </p>
        </div>

        <div className="flex flex-wrap gap-1.5">
          {sessions.map((session, index) => (
            <span
              key={index}
              className={cn(
                "inline-flex h-7 w-8 items-center justify-center rounded-md text-xs font-bold transition-colors",
                isOngoing
                  ? "bg-primary/20 text-primary"
                  : "bg-secondary text-secondary-foreground group-hover:bg-primary/20 group-hover:text-primary"
              )}
            >
              {session}
            </span>
          ))}
        </div>
      </div>
    </Link>
  );
}
