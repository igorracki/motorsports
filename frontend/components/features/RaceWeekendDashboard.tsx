"use client";

import { useEffect, useState, useCallback, useMemo } from "react";
import Link from "next/link";
import Image from "next/image";
import { ArrowLeft, Trophy, AlertCircle } from "lucide-react";
import { CircuitMap } from "@/components/features/circuit-map";
import { TrackStats } from "@/components/features/track-stats";
import { SessionSelector } from "@/components/features/session-selector";
import { ResultsTable } from "@/components/features/results-table";
import { PredictionTable } from "@/components/features/prediction-table";
import { usePredictions } from "@/hooks/usePredictions";
import { StatusBadge } from "@/components/ui/StatusBadge";
import { LoadingSpinner } from "@/components/ui/LoadingSpinner";
import { LoginModal } from "@/components/features/auth/LoginModal";
import { f1Api } from "@/services/f1-api";
import { useAuth } from "@/hooks/useAuth";
import { MainNav } from "@/components/ui/main-nav";
import { Footer } from "@/components/ui/Footer";
import { cn } from "@/lib/utils";
import type { RaceWeekend, DriverInfo, DriverResult, Circuit } from "@/types/f1";
import { getRaceStatus, isSessionLive } from "@/lib/race-utils";
import { formatRaceRange } from "@/lib/date-utils";

interface RaceWeekendDashboardProps {
  raceWeekend: RaceWeekend;
  year: number;
}

export function RaceWeekendDashboard({ raceWeekend, year }: RaceWeekendDashboardProps) {
  const [isMounted, setIsMounted] = useState(false);
  const { isAuthenticated } = useAuth();

  useEffect(() => {
    setIsMounted(true);
  }, []);
  const [showLoginModal, setShowLoginModal] = useState(false);
  const [drivers, setDrivers] = useState<DriverInfo[]>([]);
  const [circuit, setCircuit] = useState<Circuit | null>(null);
  const [sessionResults, setSessionResults] = useState<Record<string, DriverResult[]>>({});
  const [loadingResults, setLoadingResults] = useState<Record<string, boolean>>({});
  const [loadingCircuit, setLoadingCircuit] = useState(true);
  const [loadingDrivers, setLoadingDrivers] = useState(true);
  const [errorCircuit, setErrorCircuit] = useState<string | null>(null);
  const [errorDrivers, setErrorDrivers] = useState<string | null>(null);
  const [errorResults, setErrorResults] = useState<Record<string, string>>({});
  
  const fetchDashboardData = useCallback(() => {
    // Only show loading if we don't have the data yet to avoid flickering on re-fetches
    if (!drivers || drivers.length === 0) {
      setLoadingDrivers(true);
    }
    if (!circuit) {
      setLoadingCircuit(true);
    }
    setErrorDrivers(null);
    setErrorCircuit(null);

    f1Api.getDrivers(year, raceWeekend.round)
      .then(setDrivers)
      .catch(err => {
        console.error("Failed to load drivers:", err);
        setErrorDrivers("Failed to load driver entry list");
      })
      .finally(() => setLoadingDrivers(false));

    f1Api.getCircuit(year, raceWeekend.round)
      .then(setCircuit)
      .catch(err => {
        console.error("Failed to load circuit:", err);
        setErrorCircuit("Failed to load track information");
      })
      .finally(() => setLoadingCircuit(false));
  }, [year, raceWeekend.round]);

  useEffect(() => {
    fetchDashboardData();
  }, [fetchDashboardData]);

  const {
    isPredictionMode,
    selectedSession,
    currentPredictions,
    savedPredictions,
    hasChanges,
    isSubmitting,
    handleSessionSelect,
    togglePredictionMode,
    saveCurrentPredictions,
    updatePredictions
  } = usePredictions(drivers, year, raceWeekend.round);

  const fetchSessionResults = useCallback((sessionCode: string, sessionType: string) => {
    setLoadingResults(prev => ({ ...prev, [sessionType]: true }));
    setErrorResults(prev => ({ ...prev, [sessionType]: "" }));
    
    f1Api.getSessionResults(year, raceWeekend.round, sessionCode)
      .then(results => {
        setSessionResults(prev => ({ ...prev, [sessionType]: results }));
      })
      .catch(err => {
        console.error(`Failed to fetch results for ${sessionType}:`, err);
        setErrorResults(prev => ({ ...prev, [sessionType]: "Failed to load session results" }));
      })
      .finally(() => {
        setLoadingResults(prev => ({ ...prev, [sessionType]: false }));
      });
  }, [year, raceWeekend.round]);

  // Fetch all results for passed sessions on mount
  useEffect(() => {
    if (raceWeekend && !isPredictionMode) {
      const now = Date.now();
      const passedSessions = raceWeekend.sessions.filter(s => s.timeUTCMS < now);
      
      passedSessions.forEach(session => {
        const sessionCode = session.sessionCode || session.type;
        // Fetch logic wrapped in status check (ignoring current sessionResults/errorResults as dependencies to avoid loops)
        // We use setSessionResults/setErrorResults functional updates inside the promise
        setSessionResults(current => {
          if (!current[session.type]) {
             // Use internal record-based checking to avoid triggering effect again
             // However, for mount-only (mostly) background fetch, we just trigger if missing
             fetchSessionResults(sessionCode, session.type);
          }
          return current;
        });
      });
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [raceWeekend, isPredictionMode, year, raceWeekend.round, fetchSessionResults]);

  // Auto-select last passed session for completed events
  useEffect(() => {
    if (raceWeekend && !selectedSession && !isPredictionMode) {
      const now = Date.now();
      const lastPassedSession = [...raceWeekend.sessions]
        .reverse()
        .find(s => s.timeUTCMS < now);
      
      if (lastPassedSession) {
        handleSessionSelect(lastPassedSession.type);
      }
    }
  }, [raceWeekend, selectedSession, isPredictionMode, handleSessionSelect]);

  // Fetch results when session changes
  useEffect(() => {
    if (selectedSession && !isPredictionMode) {
      const session = raceWeekend.sessions.find(s => s.sessionCode === selectedSession || s.type === selectedSession);
      const sessionCode = session?.sessionCode || selectedSession;
      
      // Use setSessionResults to check existence without depending on it
      setSessionResults(current => {
        if (!current[selectedSession]) {
          // Check loading state too
          setLoadingResults(loading => {
            if (!loading[selectedSession]) {
               fetchSessionResults(sessionCode, selectedSession);
            }
            return loading;
          });
        }
        return current;
      });
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [selectedSession, isPredictionMode, raceWeekend.sessions, fetchSessionResults]);

  // Update predictions when drivers are finally loaded if they were empty
  useEffect(() => {
    if (drivers.length > 0 && currentPredictions.length === 0) {
      updatePredictions(drivers);
    }
  }, [drivers, currentPredictions.length, updatePredictions]);

  // Auto-select first upcoming session when entering prediction mode if nothing selected
  useEffect(() => {
    if (isPredictionMode && !selectedSession && raceWeekend.sessions.length > 0) {
      const now = Date.now();
      // Find first session that hasn't started yet, or the last session if all started
      const nextSession = raceWeekend.sessions.find(s => s.timeUTCMS > now) || 
                         raceWeekend.sessions[raceWeekend.sessions.length - 1];
      
      if (nextSession) {
        handleSessionSelect(nextSession.type);
      }
    }
  }, [isPredictionMode, selectedSession, raceWeekend.sessions, handleSessionSelect]);

  const handlePredictClick = () => {
    if (!isAuthenticated) {
      setShowLoginModal(true);
      return;
    }
    togglePredictionMode();
  };

  const status = isMounted
    ? getRaceStatus(year, raceWeekend.round, raceWeekend)
    : "upcoming";
  const isWeekendCompleted = status === "completed";
  
  const selectedSessionData = raceWeekend.sessions.find(s => s.sessionCode === selectedSession || s.type === selectedSession);
  const isSelectedSessionLive =
    isMounted && selectedSessionData
      ? isSessionLive(selectedSessionData.timeUTCMS)
      : false;
  const currentResults = selectedSession ? sessionResults[selectedSession] : undefined;
  
  const isSessionLocked = useMemo(() => {
    if (!isMounted || !selectedSessionData) return false;
    return selectedSessionData.timeUTCMS < Date.now();
  }, [isMounted, selectedSessionData]);
  const hasSavedPrediction = selectedSession ? !!savedPredictions[selectedSession] : false;

  return (
    <div className="min-h-screen bg-background text-foreground flex flex-col">
      <MainNav />
      <header className="sticky top-14 z-50 border-b border-border/50 bg-background/95 backdrop-blur-sm">
        <div className="mx-auto flex max-w-7xl items-center gap-4 px-4 py-4 sm:px-6 lg:px-8">
          <Link
            href={`/calendar/${year}`}
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

      <main className="mx-auto max-w-7xl px-4 py-8 sm:px-6 lg:px-8 flex-1 w-full">
        <section className="mb-8">
          <div className="grid gap-8 lg:grid-cols-2 items-stretch">
            <div className="flex flex-col">
              <div className={cn(
                "relative overflow-hidden rounded-2xl border bg-card p-6 h-full flex items-center justify-center min-h-[300px]",
                status === "ongoing" ? "border-primary/60" : "border-border/50"
              )}>
                {loadingCircuit ? (
                  <LoadingSpinner label="Loading track map..." />
                ) : errorCircuit ? (
                  <div className="flex flex-col items-center gap-2 text-muted-foreground">
                    <AlertCircle className="h-8 w-8 opacity-20" />
                    <p className="text-sm">{errorCircuit}</p>
                    <button 
                      onClick={fetchDashboardData}
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
                {loadingCircuit ? (
                  <div className="rounded-xl border border-border/50 bg-card/50 p-8 h-full flex items-center justify-center">
                    <LoadingSpinner size="sm" label="Loading track stats..." />
                  </div>
                ) : errorCircuit ? (
                  <div className="rounded-xl border border-border/50 bg-card/50 p-8 h-full flex items-center justify-center text-center">
                    <p className="text-sm text-muted-foreground">{errorCircuit}</p>
                  </div>
                ) : (
                  circuit && <TrackStats stats={circuit} />
                )}
              </div>
            </div>
          </div>
        </section>

        <section className="mb-6">
          <button
            onClick={handlePredictClick}
            type="button"
            className={cn(
              "flex w-full items-center justify-center gap-2 rounded-xl border px-6 py-4 text-base font-semibold transition-all duration-200",
              isPredictionMode
                ? "border-success bg-success text-success-foreground hover:bg-success/90"
                : "border-primary/50 bg-primary/10 text-primary hover:bg-primary/20"
            )}
          >
            <Trophy className="h-5 w-5" />
            {isPredictionMode 
              ? "Exit Prediction Mode" 
              : isWeekendCompleted ? "View Predictions" : "Predict"}
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
            sessionResults={sessionResults}
          />

          {isPredictionMode && selectedSession && (
            <div className="mt-6 animate-in fade-in slide-in-from-top-2 duration-300">
              <div className="mb-4 flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
                <div>
                  <h3 className="text-lg font-semibold">
                    {selectedSession} Prediction
                  </h3>
                  <p className="text-sm text-muted-foreground">
                    {isSessionLocked 
                      ? "This session is locked and cannot be modified." 
                      : "Drag rows to reorder your prediction"}
                  </p>
                </div>
                {!isSessionLocked && (
                  <button
                    onClick={saveCurrentPredictions}
                    disabled={!hasChanges || isSubmitting}
                    type="button"
                    className={cn(
                      "flex items-center justify-center gap-2 rounded-xl border px-4 py-2 text-sm font-semibold transition-all duration-200",
                      hasChanges && !isSubmitting
                        ? "border-success bg-success text-success-foreground hover:bg-success/90"
                        : "cursor-not-allowed border-border/30 bg-card/50 text-muted-foreground opacity-50"
                    )}
                  >
                    {isSubmitting ? (
                      <LoadingSpinner size="sm" />
                    ) : (
                      <Trophy className="h-4 w-4" />
                    )}
                    {isSubmitting 
                      ? "Saving..." 
                      : hasSavedPrediction ? "Update Prediction" : "Save Prediction"}
                  </button>
                )}
              </div>
              
              {isSessionLocked && !hasSavedPrediction ? (
                <div className="rounded-xl border border-border/50 bg-card/50 p-12 text-center">
                  <AlertCircle className="mx-auto mb-3 h-8 w-8 text-muted-foreground/50" />
                  <p className="text-muted-foreground">No prediction was submitted for this session.</p>
                </div>
              ) : loadingDrivers ? (
                <div className="rounded-xl border border-border/50 bg-card/50 p-24 text-center">
                  <LoadingSpinner label="Loading driver entry list..." />
                </div>
              ) : errorDrivers ? (
                <div className="rounded-xl border border-border/50 bg-card/50 p-12 text-center">
                  <AlertCircle className="mx-auto mb-3 h-8 w-8 text-destructive/50" />
                  <p className="text-muted-foreground">{errorDrivers}</p>
                  <button 
                    onClick={fetchDashboardData}
                    className="text-xs text-primary hover:underline mt-2"
                  >
                    Try Again
                  </button>
                </div>
              ) : drivers.length > 0 ? (
                <PredictionTable
                  key={selectedSession}
                  drivers={currentPredictions}
                  onPredictionsChange={updatePredictions}
                  readOnly={isSessionLocked}
                  totalScore={selectedSession ? savedPredictions[selectedSession]?.score : undefined}
                />
              ) : (
                <div className="rounded-xl border border-border/50 bg-card/50 p-12 text-center">
                  <AlertCircle className="mx-auto mb-3 h-8 w-8 text-muted-foreground/50" />
                  <p className="text-muted-foreground">No driver information available for this session yet.</p>
                </div>
              )}
            </div>
          )}

          {isPredictionMode && !selectedSession && (
            <div className="mt-6 rounded-xl border border-accent/50 bg-accent/10 p-8 text-center">
              <Trophy className="mx-auto mb-3 h-8 w-8 text-accent" />
              <p className="font-medium">Select a session above to view or make predictions</p>
              <p className="mt-1 text-sm text-muted-foreground">
                {isWeekendCompleted 
                  ? "Select a completed session to view your submitted prediction" 
                  : "Only upcoming sessions can be predicted"}
              </p>
            </div>
          )}

          {!isPredictionMode && selectedSession && (
            <div className="mt-6 animate-in fade-in slide-in-from-top-2 duration-300">
              {isSelectedSessionLive && (
                <div className="mb-4 flex items-center gap-2 rounded-lg border border-destructive/20 bg-destructive/5 px-4 py-2 text-destructive">
                  <div className="relative flex h-2 w-2">
                    <span className="absolute inline-flex h-full w-full animate-ping rounded-full bg-destructive opacity-75"></span>
                    <span className="relative inline-flex h-2 w-2 rounded-full bg-destructive"></span>
                  </div>
                  <span className="text-sm font-bold uppercase tracking-wider">Live Session Data</span>
                  <span className="ml-auto text-[10px] opacity-70">Refreshing every minute</span>
                </div>
              )}
              {selectedSession && loadingResults[selectedSession] ? (
                <div className="flex items-center justify-center py-24 rounded-xl border border-border/50 bg-card/50">
                  <LoadingSpinner label="Fetching session results..." />
                </div>
              ) : selectedSession && errorResults[selectedSession] ? (
                <div className="rounded-xl border border-border/50 bg-card/50 p-12 text-center">
                  <AlertCircle className="mx-auto mb-3 h-8 w-8 text-destructive/50" />
                  <p className="text-muted-foreground">{errorResults[selectedSession]}</p>
                  <button 
                    onClick={() => {
                      const session = raceWeekend.sessions.find(s => s.sessionCode === selectedSession || s.type === selectedSession);
                      const sessionCode = session?.sessionCode || selectedSession;
                      fetchSessionResults(sessionCode, selectedSession);
                    }}
                    className="text-xs text-primary hover:underline mt-2"
                  >
                    Try Again
                  </button>
                </div>
              ) : currentResults && currentResults.length > 0 ? (
                <ResultsTable
                  results={currentResults}
                  sessionName={selectedSessionData?.type || selectedSession}
                />
              ) : (
                <div className="rounded-xl border border-border/50 bg-card/50 p-12 text-center">
                  <p className="text-muted-foreground">No results available for this session yet.</p>
                </div>
              )}
            </div>
          )}
        </section>
      </main>

      <Footer />

      <LoginModal 
        isOpen={showLoginModal} 
        onClose={() => setShowLoginModal(false)}
        onSuccess={() => togglePredictionMode()}
      />
    </div>
  );
}
