"use client";

import { useEffect, useState } from "react";
import Link from "next/link";
import Image from "next/image";
import { ArrowLeft, Trophy } from "lucide-react";
import { CircuitMap } from "@/components/circuit-map";
import { TrackStats } from "@/components/track-stats";
import { SessionSelector } from "@/components/session-selector";
import { ResultsTable } from "@/components/results-table";
import { PredictionTable } from "@/components/prediction-table";
import { usePredictions } from "@/hooks/usePredictions";
import { StatusBadge } from "@/components/ui/StatusBadge";
import { f1Api } from "@/services/f1-api";
import { cn } from "@/lib/utils";
import type { RaceWeekend, DriverInfo } from "@/types/f1";
import { getRaceStatus } from "@/lib/race-utils";

interface RaceWeekendDashboardProps {
  raceWeekend: RaceWeekend;
  year: number;
}

export function RaceWeekendDashboard({ raceWeekend, year }: RaceWeekendDashboardProps) {
  const [drivers, setDrivers] = useState<DriverInfo[]>([]);
  
  useEffect(() => {
    f1Api.getDrivers().then(setDrivers);
  }, []);

  const {
    isPredictionMode,
    selectedSession,
    currentPredictions,
    savedPredictions,
    hasChanges,
    handleSessionSelect,
    togglePredictionMode,
    saveCurrentPredictions,
    updatePredictions
  } = usePredictions(drivers);

  // Update predictions when drivers are finally loaded if they were empty
  useEffect(() => {
    if (drivers.length > 0 && currentPredictions.length === 0) {
      updatePredictions(drivers);
    }
  }, [drivers, currentPredictions.length, updatePredictions]);

  const status = getRaceStatus(year, raceWeekend.round);
  const canPredict = status !== "completed";
  
  const selectedSessionData = raceWeekend.sessions.find(s => s.type === selectedSession);
  const isSessionPredictable = selectedSessionData && (!selectedSessionData.results || selectedSessionData.results.length === 0);
  const hasSavedPrediction = selectedSession ? !!savedPredictions[selectedSession] : false;

  return (
    <div className="min-h-screen bg-background text-foreground">
      <header className="sticky top-0 z-50 border-b border-border/50 bg-background/95 backdrop-blur-sm">
        <div className="mx-auto flex max-w-7xl items-center gap-4 px-4 py-4 sm:px-6 lg:px-8">
          <Link
            href={`/?year=${year}`}
            className="flex items-center gap-2 text-muted-foreground transition-colors hover:text-foreground"
          >
            <ArrowLeft className="h-5 w-5" />
            <span className="hidden sm:inline">Back to Calendar</span>
          </Link>
          <div className="h-6 w-px bg-border/50" />
          <div className="flex items-center gap-3">
            {raceWeekend.countryCode && (
              <Image
                src={`https://flagcdn.com/w80/${raceWeekend.countryCode.toLowerCase()}.png`}
                alt={`${raceWeekend.country} flag`}
                width={32}
                height={20}
                className="rounded-sm object-cover shadow-sm"
                unoptimized
              />
            )}
            <h1 className="text-lg font-bold sm:text-xl">
              {raceWeekend.fullName}
            </h1>
          </div>
          <StatusBadge status={status} className="ml-auto" />
        </div>
      </header>

      <main className="mx-auto max-w-7xl px-4 py-8 sm:px-6 lg:px-8">
        <section className="mb-8">
          <div className="grid gap-8 lg:grid-cols-2">
            <div className="flex flex-col">
              <div className={cn(
                "relative overflow-hidden rounded-2xl border bg-card p-6",
                status === "ongoing" ? "border-primary/60" : "border-border/50"
              )}>
                <CircuitMap raceWeekendId={raceWeekend.name.toLowerCase().includes("bahrain") ? `bh-${year}` : `au-${year}`} className="mx-auto aspect-[4/3] max-w-md" />
              </div>
            </div>

            <div className="flex flex-col">
              <div className="mb-6">
                <h2 className="mb-2 text-2xl font-bold sm:text-3xl">
                  {raceWeekend.name}
                </h2>
                <p className="text-muted-foreground">
                  {raceWeekend.location}, {raceWeekend.country}
                </p>
                <p className={cn(
                  "mt-2 text-lg font-medium",
                  status === "ongoing" ? "text-primary" : "text-accent"
                )}>
                  {raceWeekend.startDate || "TBC"}, {year}
                </p>
              </div>
              <TrackStats stats={{
                circuitName: raceWeekend.name,
                location: raceWeekend.location,
                country: raceWeekend.country,
                eventName: raceWeekend.fullName,
                eventDateMS: raceWeekend.startDateMS,
                lengthKM: 5.412,
                corners: 15
              }} />
            </div>
          </div>
        </section>

        <section className="mb-6">
          <button
            onClick={togglePredictionMode}
            disabled={!canPredict}
            type="button"
            className={cn(
              "flex w-full items-center justify-center gap-2 rounded-xl border px-6 py-4 text-base font-semibold transition-all duration-200",
              !canPredict
                ? "cursor-not-allowed border-border/30 bg-card/50 text-muted-foreground opacity-50"
                : isPredictionMode
                  ? "border-success bg-success text-success-foreground hover:bg-success/90"
                  : "border-primary/50 bg-primary/10 text-primary hover:bg-primary/20"
            )}
          >
            <Trophy className="h-5 w-5" />
            {isPredictionMode ? "Exit Prediction Mode" : "Predict"}
          </button>
        </section>

        <section>
          <h2 className="mb-4 text-xl font-bold">Sessions</h2>
          <SessionSelector
            sessions={raceWeekend.sessions}
            selectedSession={selectedSession}
            onSelectSession={handleSessionSelect}
            isPredictionMode={isPredictionMode}
            savedPredictions={savedPredictions}
          />

          {isPredictionMode && selectedSession && isSessionPredictable && (
            <div className="mt-6 animate-in fade-in slide-in-from-top-2 duration-300">
              <div className="mb-4 flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
                <div>
                  <h3 className="text-lg font-semibold">
                    {selectedSession} Prediction
                  </h3>
                  <p className="text-sm text-muted-foreground">
                    Drag rows to reorder your prediction
                  </p>
                </div>
                <button
                  onClick={saveCurrentPredictions}
                  disabled={!hasChanges}
                  type="button"
                  className={cn(
                    "flex items-center justify-center gap-2 rounded-xl border px-4 py-2 text-sm font-semibold transition-all duration-200",
                    hasChanges
                      ? "border-success bg-success text-success-foreground hover:bg-success/90"
                      : "cursor-not-allowed border-border/30 bg-card/50 text-muted-foreground opacity-50"
                  )}
                >
                  <Trophy className="h-4 w-4" />
                  {hasSavedPrediction ? "Update Prediction" : "Save Prediction"}
                </button>
              </div>
              <PredictionTable
                key={selectedSession}
                drivers={drivers}
                onPredictionsChange={updatePredictions}
                initialPredictions={savedPredictions[selectedSession]}
                onSave={saveCurrentPredictions}
              />
            </div>
          )}

          {isPredictionMode && !selectedSession && (
            <div className="mt-6 rounded-xl border border-accent/50 bg-accent/10 p-8 text-center">
              <Trophy className="mx-auto mb-3 h-8 w-8 text-accent" />
              <p className="font-medium">Select a session above to make your prediction</p>
              <p className="mt-1 text-sm text-muted-foreground">
                Only sessions without results can be predicted
              </p>
            </div>
          )}

          {!isPredictionMode && selectedSessionData?.results && (
            <div className="mt-6 animate-in fade-in slide-in-from-top-2 duration-300">
              <ResultsTable
                results={selectedSessionData.results}
                sessionName={selectedSessionData.type}
              />
            </div>
          )}
        </section>
      </main>
    </div>
  );
}
