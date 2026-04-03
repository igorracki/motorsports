import { Session, PredictionPolicyConfig } from "@/types/f1";

let currentConfig: PredictionPolicyConfig = {
  lockThresholdMS: 0,
  preSessionBufferMS: 15 * 60 * 1000,
  sessionDurationMS: 2 * 60 * 60 * 1000,
  revalidationWindowMS: 48 * 60 * 60 * 1000,
};

export const PredictionPolicy = {

  setConfiguration(config: PredictionPolicyConfig) {
    currentConfig = config;
  },

  getConfiguration(): PredictionPolicyConfig {
    return currentConfig;
  },

  isSessionLive(sessionTimeUTCMS: number): boolean {
    const now = Date.now();
    return (
      now >= sessionTimeUTCMS - currentConfig.preSessionBufferMS &&
      now <= sessionTimeUTCMS + currentConfig.sessionDurationMS
    );
  },

  isLocked(session: Session | { timeUTCMS: number, isLocked?: boolean }, now = Date.now()): boolean {
    if ('isLocked' in session && session.isLocked) {
      return true;
    }
    return this.hasStarted(session, now);
  },

  hasStarted(session: Session | { timeUTCMS: number }, now = Date.now()): boolean {
    return session.timeUTCMS + currentConfig.lockThresholdMS <= now;
  },

  canPredict(session: Session, now = Date.now()): boolean {
    return !this.isLocked(session, now);
  },

  isCompleted(session: Session | { timeUTCMS: number, isCompleted?: boolean }, now = Date.now()): boolean {
    if ('isCompleted' in session && session.isCompleted) {
      return true;
    }
    return now > session.timeUTCMS + currentConfig.sessionDurationMS;
  }
};
