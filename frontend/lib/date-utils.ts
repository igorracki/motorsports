const MONTHS = [
  "January", "February", "March", "April", "May", "June",
  "July", "August", "September", "October", "November", "December"
];

/**
 * Formats a date (UTC ms or local string) to the standard session time format.
 * Format: "March 14, 15:00"
 */
export function formatSessionTime(value: number | string): string {
  if (!value) return "TBC";
  
  // Per the ECMAScript specification, when a date-time string lacks an offset (no Z or +00:00),
  // it is interpreted as local time. By passing this "offset-less" string into the Date constructor,
  // we essentially trick the browser into treating the track's local hours as if they were the browser's local hours
  const date = new Date(value);
  
  return new Intl.DateTimeFormat("en-US", {
    month: "long",
    day: "numeric",
    hour: "2-digit",
    minute: "2-digit",
    hour12: false,
    timeZone: "UTC", // Standardize on UTC for the underlying date object we "tricked"
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
