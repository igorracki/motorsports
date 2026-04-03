import { RaceWeekend } from "@/types/f1";
import { PredictionPolicy } from "./policies/prediction-policy";

export type RaceStatus = "completed" | "ongoing" | "upcoming";

/**
 * Determines the status of a race weekend based on its start and end times.
 */
export function getRaceStatus(
  year: number, 
  round: number, 
  raceWeekend?: RaceWeekend,
  now = Date.now()
): RaceStatus {
  if (!raceWeekend || !raceWeekend.startDateUTCMS || !raceWeekend.endDateUTCMS) {
    // Fallback logic for when full data isn't available
    const today = new Date(now);
    if (year < today.getFullYear()) return "completed";
    if (year > today.getFullYear()) return "upcoming";
    return "upcoming";
  }

  if (now < raceWeekend.startDateUTCMS - PredictionPolicy.PRE_SESSION_BUFFER_MS) {
    return "upcoming";
  }
  
  if (now > raceWeekend.endDateUTCMS + PredictionPolicy.SESSION_DURATION_MS) {
    return "completed";
  }

  return "ongoing";
}

/**
 * Checks if a session is currently live (within a reasonable window of its start time).
 */
export function isSessionLive(sessionTimeUTCMS: number): boolean {
  return PredictionPolicy.isSessionLive(sessionTimeUTCMS);
}

/**
 * Calculates summary stats for a schedule
 */
export function getScheduleStats(raceWeekends: RaceWeekend[], year: number, now = Date.now()) {
  const stats = {
    total: raceWeekends.length,
    completed: 0,
    ongoing: 0,
    upcoming: 0,
  };

  raceWeekends.forEach((r) => {
    const status = getRaceStatus(year, r.round, r, now);
    stats[status]++;
  });

  return stats;
}
