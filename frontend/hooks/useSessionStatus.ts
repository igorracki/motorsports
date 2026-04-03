import { useState, useEffect } from "react";
import { PredictionPolicy } from "@/lib/policies/prediction-policy";
import type { Session } from "@/types/f1";

/**
 * Hook to manage session status (locked, live) while safely handling Next.js hydration.
 */
export function useSessionStatus(sessionData: Session | undefined) {
  const [isLocked, setIsLocked] = useState(false);
  const [isLive, setIsLive] = useState(false);

  useEffect(() => {
    let rAFHandle: number;
    if (sessionData) {
      rAFHandle = requestAnimationFrame(() => {
        setIsLocked(PredictionPolicy.isLocked(sessionData));
        setIsLive(PredictionPolicy.isSessionLive(sessionData.timeUTCMS));
      });
    } else {
      rAFHandle = requestAnimationFrame(() => {
        setIsLocked(false);
        setIsLive(false);
      });
    }
    return () => cancelAnimationFrame(rAFHandle);
  }, [sessionData]);

  return { isLocked, isLive };
}
