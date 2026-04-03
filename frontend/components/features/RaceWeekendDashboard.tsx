"use client";

import { useEffect, useState, useMemo } from "react";
import { Trophy } from "lucide-react";
import { SessionSelector } from "@/components/features/session-selector";
import { usePredictions } from "@/hooks/usePredictions";
import { useAuth } from "@/hooks/useAuth";
import { useModal } from "@/components/providers/modal-provider";
import { MainNav } from "@/components/ui/main-nav";
import { Footer } from "@/components/ui/Footer";
import { cn } from "@/lib/utils";
import type { RaceWeekend, DriverInfo, Circuit } from "@/types/f1";
import { getRaceStatus } from "@/lib/race-utils";
import { useRaceDashboardController } from "@/hooks/useRaceDashboardController";
import { useSessionStatus } from "@/hooks/useSessionStatus";
import { ErrorBoundary } from "@/components/ui/error-boundary";

import { DashboardHeader } from "./dashboard/DashboardHeader";
import { TrackOverview } from "./dashboard/TrackOverview";
import { PredictionWorkspace } from "./dashboard/PredictionWorkspace";
import { ResultsViewer } from "./dashboard/ResultsViewer";

interface RaceWeekendDashboardProps {
  raceWeekend: RaceWeekend;
  year: number;
  initialDrivers?: DriverInfo[];
  initialCircuit?: Circuit | null;
  serverTime?: number;
}

export function RaceWeekendDashboard({
  raceWeekend,
  year,
  initialDrivers = [],
  initialCircuit = null,
  serverTime = 0
}: RaceWeekendDashboardProps) {
  const { isAuthenticated } = useAuth();
  const { openLoginModal } = useModal();
  const [now, setNow] = useState(serverTime);

  useEffect(() => {
    // Sync with local clock after mount to avoid hydration mismatch
    const handle = requestAnimationFrame(() => setNow(Date.now()));
    return () => cancelAnimationFrame(handle);
  }, []);

  const {
    drivers,
    circuit,
    sessionResults,
    selectedSession,
    isPredictionMode,
    loadingDrivers,
    loadingCircuit,
    loadingResults,
    errorDrivers,
    errorCircuit,
    errorResults,
    fetchBaseData,
    fetchSessionResults,
    fetchPassedSessions,
    setSelectedSession,
    togglePredictionMode,
  } = useRaceDashboardController(raceWeekend, year, {
    initialDrivers,
    initialCircuit,
  });

  const {
    currentPredictions,
    savedPredictions,
    hasChanges,
    isSubmitting,
    isFetching,
    saveCurrentPredictions,
    updatePredictions,
    fetchError,
    submitError
  } = usePredictions(drivers, year, raceWeekend.round, isPredictionMode, selectedSession);

  useEffect(() => {
    if (!isPredictionMode) {
      fetchPassedSessions();
    }
  }, [isPredictionMode, fetchPassedSessions]);

  useEffect(() => {
    if (drivers.length > 0 && currentPredictions.length === 0) {
      updatePredictions(drivers);
    }
  }, [drivers, currentPredictions.length, updatePredictions]);

  const handlePredictClick = () => {
    if (!isAuthenticated) {
      openLoginModal(() => togglePredictionMode());
      return;
    }
    togglePredictionMode();
  };

  const status = useMemo(() => getRaceStatus(year, raceWeekend.round, raceWeekend, now), [raceWeekend, year, now]);
  const selectedSessionData = raceWeekend.sessions.find(s => s.sessionCode === selectedSession || s.type === selectedSession);

  const { isLocked: isSessionLocked, isLive: isSelectedSessionLive } = useSessionStatus(selectedSessionData);

  return (
    <div className="min-h-screen bg-background text-foreground flex flex-col">
      <MainNav />

      <ErrorBoundary name="Dashboard Header">
        <DashboardHeader
          raceWeekend={raceWeekend}
          year={year}
          status={status}
        />
      </ErrorBoundary>

      <main className="mx-auto max-w-7xl px-4 py-8 sm:px-6 lg:px-8 flex-1 w-full">
        <ErrorBoundary name="Track Overview">
          <TrackOverview
            raceWeekend={raceWeekend}
            year={year}
            circuit={circuit}
            status={status}
            loading={loadingCircuit}
            error={errorCircuit}
            onRetry={() => fetchBaseData(true)}
          />
        </ErrorBoundary>

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
              : status === "completed" ? "View Predictions" : "Predict"}
          </button>
        </section>

        <section>
          <h2 className="mb-4 text-xl font-bold">Sessions</h2>
          <ErrorBoundary name="Session Selection">
            <SessionSelector
              sessions={raceWeekend.sessions}
              selectedSession={selectedSession}
              onSelectSession={setSelectedSession}
              isPredictionMode={isPredictionMode}
              savedPredictions={savedPredictions}
              sessionResults={sessionResults}
              currentTime={now}
            />
          </ErrorBoundary>

          {isPredictionMode ? (
            <ErrorBoundary name="Prediction Workspace">
              <PredictionWorkspace
                session={selectedSession}
                drivers={currentPredictions}
                savedPredictions={savedPredictions}
                isSubmitting={isSubmitting}
                isFetching={isFetching}
                hasChanges={hasChanges}
                isSessionLocked={isSessionLocked}
                loadingDrivers={loadingDrivers}
                errorDrivers={errorDrivers}
                fetchError={fetchError}
                submitError={submitError}
                onSave={saveCurrentPredictions}
                onUpdatePredictions={updatePredictions}
                onRetryDrivers={() => fetchBaseData(true)}
              />
            </ErrorBoundary>
          ) : (
            <ErrorBoundary name="Results Viewer">
              <ResultsViewer
                session={selectedSession}
                sessionData={selectedSessionData}
                results={selectedSession ? sessionResults[selectedSession] : undefined}
                loading={selectedSession ? loadingResults[selectedSession] : false}
                error={selectedSession ? errorResults[selectedSession] : undefined}
                isLive={isSelectedSessionLive}
                onRetry={() => {
                  if (selectedSession) {
                    const sessionCode = selectedSessionData?.sessionCode || selectedSession;
                    fetchSessionResults(sessionCode, selectedSession);
                  }
                }}
              />
            </ErrorBoundary>
          )}
        </section>
      </main>

      <Footer />
    </div>
  );
}
