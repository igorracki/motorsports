"use client";

import { EventCard } from "./event-card";
import type { Event } from "@/lib/events-data";
import { getSessionCodes } from "@/lib/events-data";

interface EventsGridProps {
  events: Event[];
  year: number;
}

export function EventsGrid({ events, year }: EventsGridProps) {
  if (events.length === 0) {
    return (
      <div className="flex min-h-[300px] items-center justify-center">
        <p className="text-muted-foreground">
          No events scheduled for this season.
        </p>
      </div>
    );
  }

  return (
    <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
      {events.map((event) => (
        <EventCard
          key={event.id}
          id={event.id}
          title={event.title}
          country={event.country}
          countryCode={event.countryCode}
          location={event.location}
          dateFrom={event.dateFrom}
          dateTo={event.dateTo}
          sessions={getSessionCodes(event.sessions)}
          status={event.status}
          year={year}
        />
      ))}
    </div>
  );
}
