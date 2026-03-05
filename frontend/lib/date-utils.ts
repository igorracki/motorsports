/**
 * Formats a UTC millisecond timestamp to the browser's local time string.
 * Format: "March 14, 15:00"
 */
export function formatSessionTime(utcMs: number): string {
  if (!utcMs) return "TBC";
  
  const date = new Date(utcMs);
  
  return new Intl.DateTimeFormat("en-US", {
    month: "long",
    day: "numeric",
    hour: "2-digit",
    minute: "2-digit",
    hour12: false,
  }).format(date);
}

/**
 * Formats a date range from start and end UTC millisecond timestamps to local time.
 * Format: "March 14 – 16"
 */
export function formatRaceRange(startMs: number, endMs: number): string {
  if (!startMs || !endMs) return "TBC";
  
  const startDate = new Date(startMs);
  const endDate = new Date(endMs);
  
  const formatter = new Intl.DateTimeFormat("en-US", {
    day: "numeric",
    month: "long",
  });

  return formatter.formatRange(startDate, endDate);
}

/**
 * Formats a UTC millisecond timestamp to the browser's local date string.
 * Format: "March 14"
 */
export function formatDayMonth(utcMs: number): string {
  if (!utcMs) return "TBC";

  const date = new Date(utcMs);

  return new Intl.DateTimeFormat("en-US", {
    day: "numeric",
    month: "long",
  }).format(date);
}
