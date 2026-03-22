import { RaceWeekend } from "@/types/f1";

export type RaceStatus = "completed" | "ongoing" | "upcoming";

const SESSION_DURATION_MS = 2 * 60 * 60 * 1000; // Assume 2 hours max
const PRE_SESSION_MS = 15 * 60 * 1000; // 15 minutes before

/**
 * Determines the status of a race weekend based on its start and end times.
 */
export function getRaceStatus(year: number, round: number, raceWeekend?: RaceWeekend): RaceStatus {
  if (!raceWeekend || !raceWeekend.startDateUTCMS || !raceWeekend.endDateUTCMS) {
    // Fallback logic for when full data isn't available
    const now = new Date();
    if (year < now.getFullYear()) return "completed";
    if (year > now.getFullYear()) return "upcoming";
    return "upcoming";
  }

  const now = Date.now();
  
  if (now < raceWeekend.startDateUTCMS - PRE_SESSION_MS) {
    return "upcoming";
  }
  
  if (now > raceWeekend.endDateUTCMS + SESSION_DURATION_MS) {
    return "completed";
  }

  return "ongoing";
}

/**
 * Checks if a session is currently live (within a reasonable window of its start time).
 */
export function isSessionLive(sessionTimeUTCMS: number): boolean {
  const now = Date.now();
  return now >= sessionTimeUTCMS - PRE_SESSION_MS && now <= sessionTimeUTCMS + SESSION_DURATION_MS;
}

/**
 * Calculates summary stats for a schedule
 */
export function getScheduleStats(raceWeekends: RaceWeekend[], year: number) {
  const stats = {
    total: raceWeekends.length,
    completed: 0,
    ongoing: 0,
    upcoming: 0,
  };

  raceWeekends.forEach((r) => {
    const status = getRaceStatus(year, r.round, r);
    stats[status]++;
  });

  return stats;
}
