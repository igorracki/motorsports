import { useState, useCallback, useEffect, useRef } from "react";
import type { DriverInfo, Prediction } from "@/types/f1";
import { f1Api } from "@/services/f1-api";
import { useAuth } from "./useAuth";

export function usePredictions(initialDrivers: DriverInfo[], year: number, round: number) {
  const { user, isAuthenticated } = useAuth();
  const [isPredictionMode, setIsPredictionMode] = useState(false);
  const [selectedSession, setSelectedSession] = useState<string | null>(null);
  const [currentPredictions, setCurrentPredictions] = useState<DriverInfo[]>(
    initialDrivers.map(d => ({ ...d, isPredicted: false }))
  );
  // Store predictions fetched from backend: { [sessionCode]: Prediction }
  const [savedPredictions, setSavedPredictions] = useState<Record<string, Prediction>>({});
  const [isSubmitting, setIsSubmitting] = useState(false);
  
  // Track the initial state for the current session to detect changes
  const initialSessionState = useRef<string>("");

  const fetchPredictions = useCallback(async () => {
    if (!isAuthenticated || !user) return;

    try {
      const predictions = await f1Api.getRoundPredictions(user.id, year, round);
      const predictionMap: Record<string, Prediction> = {};
      predictions.forEach(p => {
        predictionMap[p.sessionType] = p;
      });
      setSavedPredictions(predictionMap);
    } catch (error) {
      console.error("Failed to fetch predictions:", error);
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
      return initialDrivers.map(d => ({ ...d, isPredicted: false }));
    }

    // Sort drivers based on saved positions
    const predictedDrivers = [...initialDrivers].map(d => ({ ...d, isPredicted: false }));
    
    // Create a map of driver_id -> position
    const positionMap = new Map(saved.entries.map(e => [e.driverId, e.position]));
    
    // Sort and mark as predicted
    return predictedDrivers.sort((a, b) => {
      const posA = positionMap.get(a.id) || 999;
      const posB = positionMap.get(b.id) || 999;
      return posA - posB;
    }).map(d => ({
      ...d,
      isPredicted: positionMap.has(d.id)
    }));
  }, [initialDrivers, savedPredictions]);

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
        .filter(d => d.isPredicted)
        .map((d, index) => ({
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

      const saved = await f1Api.submitPrediction(user.id, prediction as any);
      
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
