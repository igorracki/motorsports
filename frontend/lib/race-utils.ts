import { RaceWeekend } from "@/types/f1";

export type RaceStatus = "completed" | "ongoing" | "upcoming";

/**
 * Determines the status of a race weekend based on the year and round.
 * This is a temporary utility for dummy data to facilitate testing.
 */
export function getRaceStatus(year: number, round: number): RaceStatus {
  // Historical data (2025) is always completed
  if (year < 2026) {
    return "completed";
  }

  // Future data (2026+)
  // We'll make Round 1 "ongoing" to test the Live UI
  if (round === 1) {
    return "ongoing";
  }

  // Round 2+ are upcoming, so we can test predictions
  return "upcoming";
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
    const status = getRaceStatus(year, r.round);
    stats[status]++;
  });

  return stats;
}
