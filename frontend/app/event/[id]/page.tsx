"use client";

import { useParams, useSearchParams } from "next/navigation";
import Link from "next/link";
import { useState, useCallback, useEffect } from "react";
import { ArrowLeft, Trophy } from "lucide-react";
import { getEventById } from "@/lib/events-data";
import { CircuitMap } from "@/components/circuit-map";
import { TrackStats } from "@/components/track-stats";
import { SessionSelector } from "@/components/session-selector";
import { ResultsTable } from "@/components/results-table";
import { PredictionTable } from "@/components/prediction-table";
import { drivers, type Driver } from "@/lib/drivers-data";
import { cn } from "@/lib/utils";

export default function EventPage() {
  const params = useParams();
  const searchParams = useSearchParams();
  const id = params.id as string;
  const year = Number(searchParams.get("year")) || 2026;

  const event = getEventById(id, year);
  const [selectedSession, setSelectedSession] = useState<string | null>(null);
  const [isPredictionMode, setIsPredictionMode] = useState(false);
  const [currentPredictions, setCurrentPredictions] = useState<Driver[]>([...drivers]);
  const [savedPredictions, setSavedPredictions] = useState<Record<string, Driver[]>>({});
  const [hasChanges, setHasChanges] = useState(false);

  const canPredict = event?.status !== "completed";

  // Handle session selection - save previous predictions first (only if changed), then load new
  const handleSessionSelect = useCallback((sessionCode: string | null) => {
    // Save current session predictions before switching (if we have a current session and changes were made)
    if (selectedSession && isPredictionMode && hasChanges) {
      setSavedPredictions(prev => ({
        ...prev,
        [selectedSession]: currentPredictions
      }));
    }
    
    // Load predictions for the new session and reset change tracking
    if (sessionCode) {
      setSavedPredictions(prev => {
        const existingPredictions = prev[sessionCode];
        setCurrentPredictions(existingPredictions || [...drivers]);
        return prev;
      });
    }
    
    setHasChanges(false);
    setSelectedSession(sessionCode);
  }, [selectedSession, isPredictionMode, currentPredictions, hasChanges]);

  // Toggle prediction mode
  const handleTogglePredictionMode = useCallback(() => {
    if (isPredictionMode) {
      // Exiting prediction mode - save current predictions only if changed
      if (selectedSession && hasChanges) {
        setSavedPredictions(prev => ({
          ...prev,
          [selectedSession]: currentPredictions
        }));
      }
      setSelectedSession(null);
      setHasChanges(false);
    }
    setIsPredictionMode(!isPredictionMode);
  }, [isPredictionMode, selectedSession, currentPredictions, hasChanges]);

  // Handle explicit save button click
  const handleSavePredictions = useCallback(() => {
    if (selectedSession && hasChanges) {
      setSavedPredictions(prev => ({
        ...prev,
        [selectedSession]: currentPredictions
      }));
      setHasChanges(false);
    }
  }, [selectedSession, currentPredictions, hasChanges]);

  // Handle predictions change from table
  const handlePredictionsChange = useCallback((newPredictions: Driver[]) => {
    setCurrentPredictions(newPredictions);
    setHasChanges(true);
  }, []);

  if (!event) {
    return (
      <div className="flex min-h-screen items-center justify-center bg-background">
        <div className="text-center">
          <h1 className="mb-4 text-2xl font-bold text-foreground">
            Event not found
          </h1>
          <Link
            href="/"
            className="text-primary hover:text-primary/80 transition-colors"
          >
            Return to calendar
          </Link>
        </div>
      </div>
    );
  }

  const selectedSessionData = event.sessions.find(
    (s) => s.code === selectedSession
  );
  const isSessionPredictable = selectedSessionData && (!selectedSessionData.results || selectedSessionData.results.length === 0);
  const hasSavedPrediction = selectedSession ? !!savedPredictions[selectedSession] : false;

  return (
    <div className="min-h-screen bg-background">
      {/* Header */}
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
            <img
              src={`https://flagcdn.com/w80/${event.countryCode.toLowerCase()}.png`}
              alt={`${event.country} flag`}
              className="h-5 w-8 rounded-sm object-cover shadow-sm"
            />
            <h1 className="text-lg font-bold text-foreground sm:text-xl">
              {event.title}
            </h1>
          </div>
          {event.status === "ongoing" && (
            <div className="ml-auto flex items-center gap-1.5 rounded-md bg-primary/20 px-2 py-0.5 text-xs font-semibold text-primary">
              <span className="relative flex h-2 w-2">
                <span className="absolute inline-flex h-full w-full animate-ping rounded-full bg-primary opacity-75" />
                <span className="relative inline-flex h-2 w-2 rounded-full bg-primary" />
              </span>
              LIVE
            </div>
          )}
        </div>
      </header>

      <main className="mx-auto max-w-7xl px-4 py-8 sm:px-6 lg:px-8">
        {/* Track Info Section */}
        <section className="mb-8">
          <div className="grid gap-8 lg:grid-cols-2">
            {/* Circuit Map */}
            <div className="flex flex-col">
              <div
                className={cn(
                  "relative overflow-hidden rounded-2xl border bg-card p-6",
                  event.status === "ongoing"
                    ? "border-primary/60"
                    : "border-border/50"
                )}
              >
                <CircuitMap
                  eventId={event.id}
                  className="mx-auto aspect-[4/3] max-w-md"
                />
              </div>
            </div>

            {/* Track Details */}
            <div className="flex flex-col">
              <div className="mb-6">
                <h2 className="mb-2 text-2xl font-bold text-foreground sm:text-3xl">
                  {event.trackName}
                </h2>
                <p className="text-muted-foreground">
                  {event.location}, {event.country}
                </p>
                <p
                  className={cn(
                    "mt-2 text-lg font-medium",
                    event.status === "ongoing" ? "text-primary" : "text-accent"
                  )}
                >
                  {event.dateFrom} - {event.dateTo}, {year}
                </p>
              </div>

              <TrackStats stats={event.trackStats} />
            </div>
          </div>
        </section>

        {/* Predict Button */}
        <section className="mb-6">
          <button
            onClick={handleTogglePredictionMode}
            disabled={!canPredict}
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

        {/* Sessions Section - Full Width */}
        <section>
          <h2 className="mb-4 text-xl font-bold text-foreground">Sessions</h2>
          <SessionSelector
            sessions={event.sessions}
            selectedSession={selectedSession}
            onSelectSession={handleSessionSelect}
            isPredictionMode={isPredictionMode}
            savedPredictions={savedPredictions}
          />

          {/* Prediction Table - shown when in prediction mode and a predictable session is selected */}
          {isPredictionMode && selectedSession && isSessionPredictable && (
            <div className="mt-6 animate-in fade-in slide-in-from-top-2 duration-300">
              <div className="mb-4 flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
                <div>
                  <h3 className="text-lg font-semibold text-foreground">
                    {selectedSessionData?.name} Prediction
                  </h3>
                  <p className="text-sm text-muted-foreground">
                    Drag rows to reorder your prediction
                  </p>
                </div>
                <button
                  onClick={handleSavePredictions}
                  disabled={!hasChanges}
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
                onPredictionsChange={handlePredictionsChange}
                initialPredictions={savedPredictions[selectedSession] || [...drivers]}
              />
            </div>
          )}

          {/* Prompt to select a session when in prediction mode but no session selected */}
          {isPredictionMode && !selectedSession && (
            <div className="mt-6 rounded-xl border border-accent/50 bg-accent/10 p-8 text-center">
              <Trophy className="mx-auto mb-3 h-8 w-8 text-accent" />
              <p className="text-foreground font-medium">Select a session above to make your prediction</p>
              <p className="mt-1 text-sm text-muted-foreground">
                Only sessions without results can be predicted
              </p>
            </div>
          )}

          {/* Results Table - shown when NOT in prediction mode and session has results */}
          {!isPredictionMode && selectedSessionData?.results && (
            <div className="mt-6 animate-in fade-in slide-in-from-top-2 duration-300">
              <ResultsTable
                results={selectedSessionData.results}
                sessionName={selectedSessionData.name}
              />
            </div>
          )}
        </section>
      </main>
    </div>
  );
}
