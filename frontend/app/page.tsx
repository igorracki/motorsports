"use client";

import { useState } from "react";
import { MainNav } from "@/components/main-nav";
import { YearSelector } from "@/components/year-selector";
import { EventsGrid } from "@/components/events-grid";
import { getEventsByYear } from "@/lib/events-data";
import { Calendar, Flag, Zap } from "lucide-react";

export default function Home() {
  const [selectedYear, setSelectedYear] = useState(2026);
  const events = getEventsByYear(selectedYear);

  const totalEvents = events.length;
  const completedEvents = events.filter((e) => e.status === "completed").length;
  const ongoingEvents = events.filter((e) => e.status === "ongoing").length;

  return (
    <div className="min-h-screen bg-background">
      <MainNav />
      <YearSelector
        selectedYear={selectedYear}
        onYearChange={setSelectedYear}
      />

      <main className="container mx-auto px-4 py-8 md:px-6">
        <div className="mb-8">
          <h1 className="text-3xl font-bold tracking-tight text-foreground md:text-4xl">
            {selectedYear} Season
          </h1>
          <div className="mt-3 flex flex-wrap items-center gap-4 text-sm text-muted-foreground">
            <div className="flex items-center gap-2">
              <Calendar className="h-4 w-4 text-primary" />
              <span>{totalEvents} races scheduled</span>
            </div>
            {completedEvents > 0 && (
              <div className="flex items-center gap-2">
                <Flag className="h-4 w-4 text-accent" />
                <span>{completedEvents} races completed</span>
              </div>
            )}
            {ongoingEvents > 0 && (
              <div className="flex items-center gap-2">
                <Zap className="h-4 w-4 text-primary" />
                <span>{ongoingEvents} race in progress</span>
              </div>
            )}
          </div>
        </div>

        <EventsGrid events={events} year={selectedYear} />

        <footer className="mt-16 border-t border-border/40 pt-8 text-center">
          <p className="text-sm text-muted-foreground">
            Motorsport Calendar {selectedYear} - All dates subject to change
          </p>
        </footer>
      </main>
    </div>
  );
}
