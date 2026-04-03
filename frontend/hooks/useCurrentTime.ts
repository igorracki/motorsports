import { useState, useEffect } from "react";

/**
 * A hook that returns the current time in milliseconds.
 * Provides a deterministic fallback for SSR and hydration.
 * Allows specifying an update interval.
 *
 * @param updateIntervalMs The interval in milliseconds to update the time. Default is 0 (runs once on mount).
 * @param fallbackTime The fallback time to use before the first client-side render.
 */
export function useCurrentTime(updateIntervalMs: number = 0, fallbackTime: number = 0) {
  const [now, setNow] = useState<number>(fallbackTime);

  useEffect(() => {
    // Sync with local clock immediately on mount
    const rAFHandle = requestAnimationFrame(() => setNow(Date.now()));

    if (updateIntervalMs > 0) {
      const handle = window.setInterval(() => {
        setNow(Date.now());
      }, updateIntervalMs);
      return () => {
        clearInterval(handle);
        cancelAnimationFrame(rAFHandle);
      };
    }
    return () => cancelAnimationFrame(rAFHandle);
  }, [updateIntervalMs]);

  return now;
}
