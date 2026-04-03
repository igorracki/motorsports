import { Trophy, AlertCircle } from "lucide-react";
import { LoadingSpinner } from "@/components/ui/LoadingSpinner";
import { PredictionTable } from "@/components/features/prediction-table";
import { cn } from "@/lib/utils";
import type { DriverInfo, Prediction } from "@/types/f1";

interface PredictionWorkspaceProps {
  session: string | null;
  drivers: DriverInfo[];
  savedPredictions: Record<string, Prediction>;
  isSubmitting: boolean;
  isFetching?: boolean;
  hasChanges: boolean;
  isSessionLocked: boolean;
  loadingDrivers: boolean;
  errorDrivers: string | null;
  fetchError?: string | null;
  submitError?: string | null;
  onSave: () => Promise<void>;
  onUpdatePredictions: (drivers: DriverInfo[]) => void;
  onRetryDrivers: () => void;
}

export function PredictionWorkspace({
  session,
  drivers,
  savedPredictions,
  isSubmitting,
  isFetching,
  hasChanges,
  isSessionLocked,
  errorDrivers,
  fetchError,
  submitError,
  onSave,
  onUpdatePredictions,
  onRetryDrivers
}: PredictionWorkspaceProps) {
  if (!session) {
    return (
      <div className="mt-6 rounded-xl border border-accent/50 bg-accent/10 p-8 text-center">
        <Trophy className="mx-auto mb-3 h-8 w-8 text-accent" />
        <p className="font-medium">Select a session above to view or make predictions</p>
        <p className="mt-1 text-sm text-muted-foreground">Only upcoming sessions can be predicted</p>
      </div>
    );
  }

  const hasSavedPrediction = !!savedPredictions[session];

  return (
    <div className="mt-6 animate-in fade-in slide-in-from-top-2 duration-300">
      <div className="mb-4 flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
        <div>
          <h3 className="text-lg font-semibold">
            {session} Prediction
          </h3>
          <p className="text-sm text-muted-foreground">
            {isSessionLocked 
              ? "This session is locked and cannot be modified." 
              : "Drag rows to reorder your prediction"}
          </p>
        </div>
        {!isSessionLocked && (
          <button
            onClick={onSave}
            disabled={!hasChanges || isSubmitting || isFetching}
            type="button"
            className={cn(
              "flex items-center justify-center gap-2 rounded-xl border px-4 py-2 text-sm font-semibold transition-all duration-200",
              hasChanges && !isSubmitting && !isFetching
                ? "border-success bg-success text-success-foreground hover:bg-success/90"
                : "cursor-not-allowed border-border/30 bg-card/50 text-muted-foreground opacity-50"
            )}
          >
            {isSubmitting ? (
              <LoadingSpinner size="sm" />
            ) : (
              <Trophy className="h-4 w-4" />
            )}
            {isSubmitting 
              ? "Saving..." 
              : hasSavedPrediction ? "Update Prediction" : "Save Prediction"}
          </button>
        )}
      </div>
      
      {submitError && (
        <div className="mb-4 flex items-center gap-2 rounded-lg border border-destructive/20 bg-destructive/5 px-4 py-3 text-sm text-destructive">
          <AlertCircle className="h-4 w-4 shrink-0" />
          <p>{submitError}</p>
        </div>
      )}

      {fetchError && (
        <div className="mb-4 flex items-center gap-2 rounded-lg border border-destructive/20 bg-destructive/5 px-4 py-3 text-sm text-destructive">
          <AlertCircle className="h-4 w-4 shrink-0" />
          <p>{fetchError}</p>
        </div>
      )}

      {isFetching ? (
        <div className="rounded-xl border border-border/50 bg-card/50 p-24 text-center">
          <LoadingSpinner label="Refreshing predictions..." />
        </div>
      ) : isSessionLocked && !hasSavedPrediction ? (
        <div className="rounded-xl border border-border/50 bg-accent/5 p-12 text-center">
          <Trophy className="mx-auto mb-3 h-8 w-8 text-muted-foreground/30" />
          <p className="font-medium text-muted-foreground">No saved predictions</p>
          <p className="mt-1 text-sm text-muted-foreground/60">
            Predictions for this session are now closed.
          </p>
        </div>
      ) : errorDrivers ? (
        <div className="rounded-xl border border-border/50 bg-card/50 p-12 text-center">
          <AlertCircle className="mx-auto mb-3 h-8 w-8 text-destructive/50" />
          <p className="text-muted-foreground">{errorDrivers}</p>
          <button 
            onClick={onRetryDrivers}
            className="text-xs text-primary hover:underline mt-2"
          >
            Try Again
          </button>
        </div>
      ) : drivers.length > 0 ? (
        <PredictionTable
          key={session}
          drivers={drivers}
          onPredictionsChange={onUpdatePredictions}
          readOnly={isSessionLocked}
          totalScore={savedPredictions[session]?.score}
        />
      ) : (
        <div className="rounded-xl border border-border/50 bg-card/50 p-12 text-center">
          <AlertCircle className="mx-auto mb-3 h-8 w-8 text-muted-foreground/50" />
          <p className="text-muted-foreground">No driver information available for this session yet.</p>
        </div>
      )}
    </div>
  );
}
