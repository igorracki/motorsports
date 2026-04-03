import { DriverInfo, Prediction, SessionScoringRules } from "@/types/f1";

/**
 * Mapping session codes to scoring rule types.
 * Decouples frontend session keys from backend rule identifiers.
 */
export const PredictionSessionMapper = {
  /**
   * Matches a session code to the corresponding scoring rules.
   */
  matchRules(rules: SessionScoringRules[], sessionCode: string): SessionScoringRules | undefined {
    return rules.find(r =>
      r.sessionType === sessionCode ||
      (sessionCode === "R" && r.sessionType === "Race") ||
      (sessionCode === "S" && r.sessionType === "Sprint") ||
      (sessionCode === "Q" && r.sessionType === "Qualifying") ||
      (sessionCode === "SQ" && r.sessionType === "Sprint Qualifying") ||
      (sessionCode.startsWith("FP") && r.sessionType.startsWith("Practice"))
    );
  },

  /**
   * Merges initial driver list with saved predictions.
   */
  mapDriversWithPredictions(
    initialDrivers: DriverInfo[],
    savedPrediction: Prediction | undefined,
    sessionRules: SessionScoringRules | undefined
  ): DriverInfo[] {
    if (!savedPrediction || !savedPrediction.entries) {
      return initialDrivers.map(d => ({ ...d, isPredicted: false, correct: false, points: 0 }));
    }

    const posMap = new Map(savedPrediction.entries.map(e => [e.position, e]));
    const predictedDriverIds = new Set(savedPrediction.entries.map(e => e.driverId));
    const availableDrivers = initialDrivers.filter(d => !predictedDriverIds.has(d.id));

    const result: DriverInfo[] = [];
    let availableIdx = 0;

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
  },

  /**
   * Performs a targeted, shallow comparison of two driver arrays.
   * Decouples change detection from expensive JSON serialization.
   */
  hasPredictionsChanged(initial: DriverInfo[], current: DriverInfo[]): boolean {
    if (initial.length !== current.length) return true;
    for (let i = 0; i < initial.length; i++) {
      if (initial[i].id !== current[i].id) return true;
      if (initial[i].isPredicted !== current[i].isPredicted) return true;
    }
    return false;
  }
};
