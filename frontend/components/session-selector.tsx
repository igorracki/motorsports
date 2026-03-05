"use client";

import { cn } from "@/lib/utils";
import { Trophy, Radio } from "lucide-react";
import type { Session, DriverInfo, DriverResult } from "@/types/f1";
import { formatSessionTime } from "@/lib/date-utils";
import { isSessionLive } from "@/lib/race-utils";

interface SessionSelectorProps {
  sessions: Session[];
  selectedSession: string | null;
  onSelectSession: (code: string | null) => void;
  isPredictionMode: boolean;
  savedPredictions: Record<string, DriverInfo[]>;
  sessionResults?: Record<string, DriverResult[]>;
}

export function SessionSelector({
  sessions,
  selectedSession,
  onSelectSession,
  isPredictionMode,
  savedPredictions,
  sessionResults = {},
}: SessionSelectorProps) {
  const now = Date.now();

  return (
    <div className="grid grid-cols-5 gap-2">
      {sessions.map((session) => {
        const isSelected = selectedSession === session.type;
        const results = sessionResults[session.type] || session.results;
        const hasResults = results && results.length > 0;
        const hasPrediction = !!savedPredictions[session.type];
        
        const isStarted = session.timeUTCMS < now;
        const isLive = isSessionLive(session.timeUTCMS);
        
        // 1. Normal mode: Always clickable
        // 2. Prediction mode: Only clickable if it hasn't started yet
        const isClickable = isPredictionMode ? !isStarted : true;

        return (
          <button
            key={session.type}
            onClick={() => {
              if (isClickable) {
                onSelectSession(isSelected ? null : session.type);
              }
            }}
            type="button"
            className={cn(
              "group relative flex flex-col items-center justify-center rounded-xl border px-2 py-3 transition-all duration-200 sm:px-4",
              isSelected
                ? "border-primary bg-primary/20 text-primary shadow-sm shadow-primary/10"
                : isClickable
                  ? isPredictionMode
                    ? "border-accent/50 bg-accent/10 text-foreground hover:border-accent hover:bg-accent/20"
                    : "border-border/50 bg-card text-foreground hover:border-primary/50 hover:bg-primary/10"
                  : "cursor-not-allowed border-border/30 bg-card/50 text-muted-foreground opacity-50"
            )}
            disabled={!isClickable}
          >
            {hasPrediction && !hasResults && (
              <div className="absolute -right-1 -top-1 flex h-4 w-4 items-center justify-center rounded-full bg-accent text-accent-foreground border-2 border-background">
                <Trophy className="h-2.5 w-2.5" />
              </div>
            )}
            {isLive && !isPredictionMode && (
              <div className="absolute -left-1 -top-1 flex h-4 w-auto items-center justify-center rounded-full bg-destructive px-1.5 text-[8px] font-black uppercase tracking-tighter text-destructive-foreground border-2 border-background animate-pulse">
                Live
              </div>
            )}
            <span className="text-xs font-bold sm:text-sm">{session.sessionCode || session.type}</span>
            <span
              className={cn(
                "mt-1 text-[10px] sm:text-xs font-medium",
                isSelected ? "text-primary" : "text-muted-foreground"
              )}
            >
              {formatSessionTime(session.timeUTCMS)}
            </span>
            
            <span
              className={cn(
                "mt-1 text-[10px] uppercase tracking-wider font-bold",
                isPredictionMode
                  ? isSelected ? "text-primary" : "text-accent"
                  : isStarted
                    ? hasResults ? "text-success" : (sessionResults[session.type] ? "text-muted-foreground" : "text-success")
                    : "text-muted-foreground"
              )}
            >
              {isPredictionMode
                ? hasPrediction ? "Predicted" : "Predict"
                : isLive
                  ? "Live"
                  : isStarted 
                    ? (sessionResults[session.type] && !hasResults ? "No data" : "Results") 
                    : "Upcoming"}
            </span>
          </button>
        );
      })}
    </div>
  );
}
