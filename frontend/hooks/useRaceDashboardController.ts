import { useState, useEffect, useCallback } from "react";
import { useApi } from "@/components/providers/api-provider";
import type { RaceWeekend, DriverInfo, Circuit } from "@/types/f1";
import { PredictionPolicy } from "@/lib/policies/prediction-policy";
import { useAsync } from "./useAsync";
import { useKeyedAsync } from "./useKeyedAsync";

interface UseRaceDashboardControllerOptions {
  initialDrivers?: DriverInfo[];
  initialCircuit?: Circuit | null;
}

/**
 * Controller hook to manage the data orchestration for the Race Weekend Dashboard.
 * Separates API interaction and state management from UI rendering.
 */
export function useRaceDashboardController(
  raceWeekend: RaceWeekend,
  year: number,
  options: UseRaceDashboardControllerOptions = {}
) {
  const { raceRepo } = useApi();
  const [isPredictionMode, setIsPredictionMode] = useState(false);
  const [selectedSession, setSelectedSession] = useState<string | null>(null);

  const {
    execute: executeFetchDrivers,
    data: fetchedDrivers,
    loading: isLoadingDrivers,
    error: errorDrivers
  } = useAsync(() => raceRepo.getDrivers(year, raceWeekend.round));

  const {
    execute: executeFetchCircuit,
    data: fetchedCircuit,
    loading: isLoadingCircuit,
    error: errorCircuit
  } = useAsync(() => raceRepo.getCircuit(year, raceWeekend.round));

  const {
    execute: executeFetchSessionResults,
    data: sessionResults,
    loading: loadingResults,
    error: errorResults,
  } = useKeyedAsync(async (sessionType: string, sessionCode: string) => {
    return raceRepo.getSessionResults(year, raceWeekend.round, sessionCode);
  });

  const drivers = fetchedDrivers || options.initialDrivers || [];
  const circuit = fetchedCircuit || options.initialCircuit || null;

  const loadingDrivers = isLoadingDrivers && drivers.length === 0;
  const loadingCircuit = isLoadingCircuit && circuit === null;

  const fetchBaseData = useCallback(async (force = false) => {
    const promises = [];
    if (force || drivers.length === 0) promises.push(executeFetchDrivers());
    if (force || circuit === null) promises.push(executeFetchCircuit());
    if (promises.length > 0) {
      await Promise.allSettled(promises);
    }
  }, [drivers.length, circuit, executeFetchDrivers, executeFetchCircuit]);

  const fetchSessionResults = useCallback(async (sessionCode: string, sessionType: string) => {
    if (sessionResults[sessionType] || loadingResults[sessionType]) return;
    try {
      await executeFetchSessionResults(sessionType, sessionCode);
    } catch {
      // Error handled by useKeyedAsync
    }
  }, [sessionResults, loadingResults, executeFetchSessionResults]);

  const fetchPassedSessions = useCallback(() => {
    const now = Date.now();
    const passedSessions = raceWeekend.sessions.filter(s => PredictionPolicy.hasStarted(s, now));

    passedSessions.forEach(session => {
      const sessionCode = session.sessionCode || session.type;
      if (!sessionResults[session.type] && !loadingResults[session.type]) {
        fetchSessionResults(sessionCode, session.type);
      }
    });
  }, [raceWeekend.sessions, sessionResults, loadingResults, fetchSessionResults]);

  useEffect(() => {
    let rAFHandle: number;
    if (!selectedSession && !isPredictionMode) {
      const now = Date.now();
      const lastPassedSession = [...raceWeekend.sessions]
        .reverse()
        .find(s => PredictionPolicy.hasStarted(s, now));

      if (lastPassedSession) {
        rAFHandle = requestAnimationFrame(() => setSelectedSession(lastPassedSession.type));
      }
    } else if (isPredictionMode && !selectedSession && raceWeekend.sessions.length > 0) {
      const now = Date.now();
      const nextSession = raceWeekend.sessions.find(s => !PredictionPolicy.hasStarted(s, now)) ||
        raceWeekend.sessions[raceWeekend.sessions.length - 1];

      if (nextSession) {
        rAFHandle = requestAnimationFrame(() => setSelectedSession(nextSession.type));
      }
    }
    return () => cancelAnimationFrame(rAFHandle);
  }, [raceWeekend, selectedSession, isPredictionMode]);

  const togglePredictionMode = useCallback(() => {
    setIsPredictionMode(prev => {
      const nextMode = !prev;
      if (nextMode) setSelectedSession(null);
      return nextMode;
    });
  }, []);

  // Sync session results when selection changes
  useEffect(() => {
    if (selectedSession && !isPredictionMode) {
      const session = raceWeekend.sessions.find(s => s.sessionCode === selectedSession || s.type === selectedSession);
      const sessionCode = session?.sessionCode || selectedSession;

      fetchSessionResults(sessionCode, selectedSession);
    }
  }, [selectedSession, isPredictionMode, raceWeekend.sessions, fetchSessionResults]);

  useEffect(() => {
    fetchBaseData();
  }, [fetchBaseData]);

  return {
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
  };
}
