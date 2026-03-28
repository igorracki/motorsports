import { useState, useCallback, useEffect, useRef } from "react";
import type { DriverInfo, Prediction, SessionScoringRules } from "@/types/f1";
import { f1Api } from "@/services/f1-api";
import { useAuth } from "./useAuth";

export function usePredictions(initialDrivers: DriverInfo[], year: number, round: number) {
  const { user, isAuthenticated } = useAuth();
  const [isPredictionMode, setIsPredictionMode] = useState(false);
  const [selectedSession, setSelectedSession] = useState<string | null>(null);
  const [scoringRules, setScoringRules] = useState<SessionScoringRules[]>([]);
  const [currentPredictions, setCurrentPredictions] = useState<DriverInfo[]>(
    initialDrivers.map(d => ({ ...d, isPredicted: false, correct: false, points: 0 }))
  );
  // Store predictions fetched from backend: { [sessionCode]: Prediction }
  const [savedPredictions, setSavedPredictions] = useState<Record<string, Prediction>>({});
  const [isSubmitting, setIsSubmitting] = useState(false);
  
  // Track the initial state for the current session to detect changes
  const initialSessionState = useRef<string>("");

  const fetchPredictions = useCallback(async () => {
    if (!isAuthenticated || !user) return;

    try {
      const [predictions, rules] = await Promise.all([
        f1Api.getRoundPredictions(user.id, year, round),
        f1Api.getScoringRules()
      ]);
      
      const predictionMap: Record<string, Prediction> = {};
      predictions.forEach(p => {
        predictionMap[p.sessionType] = p;
      });
      setSavedPredictions(predictionMap);
      setScoringRules(rules);
    } catch (error) {
      console.error("Failed to fetch predictions or rules:", error);
    }
  }, [isAuthenticated, user, year, round]);

  // Fetch predictions when entering prediction mode
  useEffect(() => {
    if (isPredictionMode) {
      fetchPredictions();
    }
  }, [isPredictionMode, fetchPredictions]);

  const getDriverListWithPredictions = useCallback((sessionCode: string) => {
    const saved = savedPredictions[sessionCode];
    if (!saved || !saved.entries) {
      return initialDrivers.map(d => ({ ...d, isPredicted: false, correct: false, points: 0 }));
    }

    const posMap = new Map(saved.entries.map(e => [e.position, e]));
    const predictedDriverIds = new Set(saved.entries.map(e => e.driverId));
    const availableDrivers = initialDrivers.filter(d => !predictedDriverIds.has(d.id));
    
    const sessionRules = scoringRules.find(r => 
      r.sessionType === sessionCode || 
      (sessionCode === "R" && r.sessionType === "Race") ||
      (sessionCode === "S" && r.sessionType === "Sprint") ||
      (sessionCode === "Q" && r.sessionType === "Qualifying") ||
      (sessionCode === "SQ" && r.sessionType === "Sprint Qualifying") ||
      (sessionCode.startsWith("FP") && r.sessionType.startsWith("Practice"))
    );

    const result: DriverInfo[] = [];
    let availableIdx = 0;

    // We assume the total number of drivers is initialDrivers.length
    for (let i = 1; i <= initialDrivers.length; i++) {
      const pred = posMap.get(i);
      if (pred) {
        const driver = initialDrivers.find(d => d.id === pred.driverId);
        if (driver) {
          let points = 0;
          if (pred.correct && sessionRules) {
            const rule = sessionRules.rules.find(r => r.position === pred.position);
            points = rule ? rule.points : 0;
          }
          result.push({
            ...driver,
            isPredicted: true,
            correct: pred.correct || false,
            points
          });
          continue;
        }
      }
      
      // No prediction for this position, take the next available driver
      if (availableIdx < availableDrivers.length) {
        result.push({
          ...availableDrivers[availableIdx],
          isPredicted: false,
          correct: false,
          points: 0
        });
        availableIdx++;
      }
    }

    return result;
  }, [initialDrivers, savedPredictions, scoringRules]);

  // Sync currentPredictions when savedPredictions are loaded or session changes while in prediction mode
  useEffect(() => {
    if (isPredictionMode && selectedSession) {
      const drivers = getDriverListWithPredictions(selectedSession);
      setCurrentPredictions(drivers);
      initialSessionState.current = JSON.stringify(drivers);
    }
  }, [isPredictionMode, selectedSession, savedPredictions, getDriverListWithPredictions]);

  const handleSessionSelect = useCallback((sessionCode: string | null) => {
    if (sessionCode) {
      const drivers = getDriverListWithPredictions(sessionCode);
      setCurrentPredictions(drivers);
      initialSessionState.current = JSON.stringify(drivers);
    }
    setSelectedSession(sessionCode);
  }, [getDriverListWithPredictions]);

  const togglePredictionMode = useCallback(() => {
    if (isPredictionMode) {
      setSelectedSession(null);
    }
    setIsPredictionMode(prev => !prev);
  }, [isPredictionMode]);

  const hasChanges = JSON.stringify(currentPredictions) !== initialSessionState.current;

  const saveCurrentPredictions = useCallback(async () => {
    if (!selectedSession || !user || !hasChanges || isSubmitting) return;

    setIsSubmitting(true);
    try {
      const entries = currentPredictions
        .map((d, index) => ({ d, index }))
        .filter(item => item.d.isPredicted)
        .map(({ d, index }) => ({
          prediction_id: "", // Will be set by backend
          position: index + 1,
          driver_id: d.id,
        }));

      const prediction = {
        year,
        round,
        session_type: selectedSession,
        entries: entries,
      };

      const saved = await f1Api.submitPrediction(user.id, prediction);
      
      setSavedPredictions(prev => ({
        ...prev,
        [selectedSession]: saved
      }));
      
      initialSessionState.current = JSON.stringify(currentPredictions);
    } catch (error) {
      console.error("Failed to save prediction:", error);
    } finally {
      setIsSubmitting(false);
    }
  }, [selectedSession, user, currentPredictions, hasChanges, isSubmitting, year, round]);

  const updatePredictions = useCallback((newPredictions: DriverInfo[]) => {
    setCurrentPredictions(newPredictions);
  }, []);

  return {
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
  };
}
