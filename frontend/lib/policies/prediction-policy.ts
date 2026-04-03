import { Session } from "@/types/f1";

export const PredictionPolicy = {
  /**
   * Defines the window before a session starts during which it is considered "live".
   */
  PRE_SESSION_BUFFER_MS: 15 * 60 * 1000, // 15 minutes

  /**
   * Defines the assumed maximum duration of a session.
   */
  SESSION_DURATION_MS: 2 * 60 * 60 * 1000, // 2 hours

  /**
   * Determines if a session is currently "live" (ongoing or about to start).
   */
  isSessionLive(sessionTimeUTCMS: number): boolean {
    const now = Date.now();
    return (
      now >= sessionTimeUTCMS - this.PRE_SESSION_BUFFER_MS &&
      now <= sessionTimeUTCMS + this.SESSION_DURATION_MS
    );
  },

  /**
   * Determines if a session is "locked" for predictions.
   * Rules:
   * 1. A session is locked once it has started.
   */
  isLocked(session: Session | { timeUTCMS: number }, now = Date.now()): boolean {
    return this.hasStarted(session, now);
  },

  /**
   * Determines if a session has already started or passed.
   */
  hasStarted(session: Session | { timeUTCMS: number }, now = Date.now()): boolean {
    return session.timeUTCMS < now;
  },

  /**
   * Determines if a user is allowed to submit or modify a prediction for a session.
   */
  canPredict(session: Session, now = Date.now()): boolean {
    return !this.isLocked(session, now);
  }
};
