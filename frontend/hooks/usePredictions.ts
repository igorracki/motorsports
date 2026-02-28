import { useState, useCallback } from "react";
import type { DriverInfo } from "@/types/f1";

export function usePredictions(initialDrivers: DriverInfo[]) {
  const [isPredictionMode, setIsPredictionMode] = useState(false);
  const [selectedSession, setSelectedSession] = useState<string | null>(null);
  const [currentPredictions, setCurrentPredictions] = useState<DriverInfo[]>([...initialDrivers]);
  const [savedPredictions, setSavedPredictions] = useState<Record<string, DriverInfo[]>>({});
  const [hasChanges, setHasChanges] = useState(false);

  const saveToStore = useCallback((sessionCode: string, predictions: DriverInfo[]) => {
    setSavedPredictions(prev => ({
      ...prev,
      [sessionCode]: predictions
    }));
  }, []);

  const handleSessionSelect = useCallback((sessionCode: string | null) => {
    // Save current session predictions before switching (if changes were made)
    if (selectedSession && isPredictionMode && hasChanges) {
      saveToStore(selectedSession, currentPredictions);
    }
    
    // Load predictions for the new session
    if (sessionCode) {
      const existing = savedPredictions[sessionCode];
      setCurrentPredictions(existing || [...initialDrivers]);
    }
    
    setHasChanges(false);
    setSelectedSession(sessionCode);
  }, [selectedSession, isPredictionMode, currentPredictions, hasChanges, savedPredictions, initialDrivers, saveToStore]);

  const togglePredictionMode = useCallback(() => {
    if (isPredictionMode) {
      // Exiting - save if needed
      if (selectedSession && hasChanges) {
        saveToStore(selectedSession, currentPredictions);
      }
      setSelectedSession(null);
      setHasChanges(false);
    }
    setIsPredictionMode(prev => !prev);
  }, [isPredictionMode, selectedSession, currentPredictions, hasChanges, saveToStore]);

  const saveCurrentPredictions = useCallback(() => {
    if (selectedSession && hasChanges) {
      saveToStore(selectedSession, currentPredictions);
      setHasChanges(false);
    }
  }, [selectedSession, currentPredictions, hasChanges, saveToStore]);

  const updatePredictions = useCallback((newPredictions: DriverInfo[]) => {
    setCurrentPredictions(newPredictions);
    setHasChanges(true);
  }, []);

  return {
    isPredictionMode,
    selectedSession,
    currentPredictions,
    savedPredictions,
    hasChanges,
    handleSessionSelect,
    togglePredictionMode,
    saveCurrentPredictions,
    updatePredictions
  };
}
