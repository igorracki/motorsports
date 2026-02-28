"use client";

import { cn } from "@/lib/utils";
import { Trophy } from "lucide-react";
import type { Session, DriverInfo } from "@/types/f1";

interface SessionSelectorProps {
  sessions: Session[];
  selectedSession: string | null;
  onSelectSession: (code: string | null) => void;
  isPredictionMode: boolean;
  savedPredictions: Record<string, DriverInfo[]>;
}

export function SessionSelector({
  sessions,
  selectedSession,
  onSelectSession,
  isPredictionMode,
  savedPredictions,
}: SessionSelectorProps) {
  return (
    <div className="grid grid-cols-5 gap-2">
      {sessions.map((session) => {
        const isSelected = selectedSession === session.type;
        const hasResults = session.results && session.results.length > 0;
        const hasPrediction = !!savedPredictions[session.type];
        const canSelectInPredictionMode = isPredictionMode && !hasResults;
        const canSelectInNormalMode = !isPredictionMode && hasResults;
        const isClickable = canSelectInPredictionMode || canSelectInNormalMode;

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
                ? "border-primary bg-primary/20 text-primary"
                : isClickable
                  ? isPredictionMode
                    ? "border-accent/50 bg-accent/10 text-foreground hover:border-accent hover:bg-accent/20"
                    : "border-border/50 bg-card text-foreground hover:border-primary/50 hover:bg-primary/10"
                  : "cursor-not-allowed border-border/30 bg-card/50 text-muted-foreground opacity-50"
            )}
            disabled={!isClickable}
          >
            {hasPrediction && !hasResults && (
              <div className="absolute -right-1 -top-1 flex h-4 w-4 items-center justify-center rounded-full bg-accent text-accent-foreground">
                <Trophy className="h-2.5 w-2.5" />
              </div>
            )}
            <span className="text-xs font-bold sm:text-sm">{session.type}</span>
            <span
              className={cn(
                "mt-1 text-[10px] sm:text-xs",
                isSelected ? "text-primary" : "text-muted-foreground"
              )}
            >
              {session.timeLocal || "TBC"}
            </span>
            {!hasResults && (
              <span
                className={cn(
                  "mt-1 text-[10px] uppercase tracking-wider",
                  isPredictionMode
                    ? isSelected
                      ? "text-primary"
                      : "text-accent"
                    : "text-muted-foreground"
                )}
              >
                {isPredictionMode
                  ? hasPrediction
                    ? "Predicted"
                    : "Predict"
                  : "No data"}
              </span>
            )}
          </button>
        );
      })}
    </div>
  );
}
