import { RaceWeekend } from "@/types/f1";

export type RaceStatus = "completed" | "ongoing" | "upcoming";

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
  
  if (now < raceWeekend.startDateUTCMS) {
    return "upcoming";
  }
  
  if (now > raceWeekend.endDateUTCMS) {
    return "completed";
  }

  return "ongoing";
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
