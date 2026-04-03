import { useState, useCallback, useEffect } from "react";
import type { DriverInfo, Prediction, SessionScoringRules, SubmitPredictionRequest } from "@/types/f1";
import { useApi } from "@/components/providers/api-provider";
import { useConfig } from "@/components/providers/config-provider";
import { useAuth } from "./useAuth";
import { PredictionSessionMapper } from "@/lib/mappers/prediction-mapper";
import { useAsync } from "./useAsync";

/**
 * Hook to manage prediction state and business logic for a specific race round.
 */
export function usePredictions(
  initialDrivers: DriverInfo[], 
  year: number, 
  round: number,
  isPredictionMode: boolean,
  selectedSession: string | null
) {
  const { predictionRepo } = useApi();
  const { config } = useConfig();
  const { user, isAuthenticated } = useAuth();
  const [scoringRules, setScoringRules] = useState<SessionScoringRules[]>([]);
  const [currentPredictions, setCurrentPredictions] = useState<DriverInfo[]>(
    initialDrivers.map(d => ({ ...d, isPredicted: false, correct: false, points: 0 }))
  );
  
  const [savedPredictions, setSavedPredictions] = useState<Record<string, Prediction>>({});
  const [initialSessionState, setInitialSessionState] = useState<DriverInfo[]>([]);

  const { execute: fetchAllData, loading: isFetching, error: fetchError } = useAsync(
    async () => {
      if (!user) return;
      return Promise.all([
        predictionRepo.getRoundPredictions(user.id, year, round),
        predictionRepo.getScoringRules()
      ]);
    },
    {
      onSuccess: (data) => {
        if (!data) return;
        const [predictions, rules] = data;
        const predictionMap: Record<string, Prediction> = {};
        predictions.forEach(p => {
          predictionMap[p.sessionType] = p;
        });
        setSavedPredictions(predictionMap);
        setScoringRules(rules);
      }
    }
  );

  const { execute: submitPrediction, loading: isSubmitting, error: submitError } = useAsync(
    async (predictionData: SubmitPredictionRequest) => {
      if (!user) return;
      return predictionRepo.submitPrediction(user.id, predictionData);
    },
    {
      onSuccess: (saved) => {
        if (saved && selectedSession) {
          setSavedPredictions(prev => ({
            ...prev,
            [selectedSession]: saved
          }));
          setInitialSessionState(currentPredictions);
        }
      }
    }
  );

  useEffect(() => {
    if (isPredictionMode && isAuthenticated) {
      fetchAllData();
    }
  }, [isPredictionMode, isAuthenticated, fetchAllData]);

  const getDriverListWithPredictions = useCallback((sessionCode: string) => {
    const saved = savedPredictions[sessionCode];
    const sessionRules = PredictionSessionMapper.matchRules(scoringRules, sessionCode, config?.sessionMappings);
    
    return PredictionSessionMapper.mapDriversWithPredictions(
      initialDrivers,
      saved,
      sessionRules
    );
  }, [initialDrivers, savedPredictions, scoringRules, config?.sessionMappings]);

  // Sync currentPredictions when savedPredictions are loaded or session changes while in prediction mode
  useEffect(() => {
    let rAFHandle: number;
    if (isPredictionMode && selectedSession) {
      rAFHandle = requestAnimationFrame(() => {
        const drivers = getDriverListWithPredictions(selectedSession);
        setCurrentPredictions(drivers);
        setInitialSessionState(drivers);
      });
    }
    return () => cancelAnimationFrame(rAFHandle);
  }, [isPredictionMode, selectedSession, savedPredictions, getDriverListWithPredictions]);

  const hasChanges = PredictionSessionMapper.hasPredictionsChanged(
    initialSessionState, 
    currentPredictions
  );

  const saveCurrentPredictions = useCallback(async () => {
    if (!selectedSession || !user || !hasChanges || isSubmitting) return;

    const entries = currentPredictions
      .map((d, index) => ({ d, index }))
      .filter(item => item.d.isPredicted)
      .map(({ d, index }) => ({
        prediction_id: "", 
        position: index + 1,
        driver_id: d.id,
      }));

    const prediction = {
      year,
      round,
      session_type: selectedSession,
      entries: entries,
    };

    await submitPrediction(prediction);
  }, [selectedSession, user, currentPredictions, hasChanges, isSubmitting, year, round, submitPrediction]);

  const updatePredictions = useCallback((newPredictions: DriverInfo[]) => {
    setCurrentPredictions(newPredictions);
  }, []);

  return {
    currentPredictions,
    savedPredictions,
    hasChanges,
    isSubmitting,
    isFetching,
    saveCurrentPredictions,
    updatePredictions,
    fetchError,
    submitError
  };
}
