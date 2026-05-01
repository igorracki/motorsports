const MONTHS = [
  "January", "February", "March", "April", "May", "June",
  "July", "August", "September", "October", "November", "December"
];

/**
 * Formats a date (UTC ms) to the standard session time format.
 * Format: "March 14, 15:00"
 */
export function formatSessionTime(
  utcMs: number,
  timeZone: string = "UTC"
): string {
  if (!utcMs) return "TBC";
  
  const date = new Date(utcMs);
  
  return new Intl.DateTimeFormat("en-US", {
    month: "long",
    day: "numeric",
    hour: "2-digit",
    minute: "2-digit",
    hour12: false,
    timeZone: timeZone,
  }).format(date);
}

/**
 * Formats a date range from start and end UTC millisecond timestamps.
 * This version is deterministic to prevent hydration mismatches.
 * Format: "March 14 - 16"
 */
export function formatRaceRange(startMs: number, endMs: number): string {
  if (!startMs || !endMs) return "TBC";
  
  const startDate = new Date(startMs);
  const endDate = new Date(endMs);
  
  const startMonth = MONTHS[startDate.getUTCMonth()];
  const startDay = startDate.getUTCDate();
  const endMonth = MONTHS[endDate.getUTCMonth()];
  const endDay = endDate.getUTCDate();

  if (startMonth === endMonth) {
    return `${startMonth} ${startDay} - ${endDay}`;
  }
  
  return `${startMonth} ${startDay} - ${endMonth} ${endDay}`;
}

/**
 * Formats a UTC millisecond timestamp to a deterministic date string.
 * Format: "March 14"
 */
export function formatDayMonth(utcMs: number): string {
  if (!utcMs) return "TBC";

  const date = new Date(utcMs);
  const month = MONTHS[date.getUTCMonth()];
  const day = date.getUTCDate();

  return `${month} ${day}`;
}
