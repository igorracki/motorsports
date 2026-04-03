import { AlertCircle } from "lucide-react";
import { LoadingSpinner } from "@/components/ui/LoadingSpinner";
import { ResultsTable } from "@/components/features/results-table";
import type { DriverResult, Session } from "@/types/f1";

interface ResultsViewerProps {
  session: string | null;
  sessionData?: Session;
  results?: DriverResult[];
  loading: boolean;
  error?: string;
  isLive: boolean;
  onRetry: () => void;
}

export function ResultsViewer({
  session,
  sessionData,
  results,
  loading,
  error,
  isLive,
  onRetry
}: ResultsViewerProps) {
  if (!session) {
    return (
      <div className="mt-6 rounded-xl border border-border/50 bg-card/50 p-12 text-center">
        <p className="text-muted-foreground">Select a session above to view results.</p>
      </div>
    );
  }

  return (
    <div className="mt-6 animate-in fade-in slide-in-from-top-2 duration-300">
      {isLive && (
        <div className="mb-4 flex items-center gap-2 rounded-lg border border-destructive/20 bg-destructive/5 px-4 py-2 text-destructive">
          <div className="relative flex h-2 w-2">
            <span className="absolute inline-flex h-full w-full animate-ping rounded-full bg-destructive opacity-75"></span>
            <span className="relative inline-flex h-2 w-2 rounded-full bg-destructive"></span>
          </div>
          <span className="text-sm font-bold uppercase tracking-wider">Live Session Data</span>
          <span className="ml-auto text-[10px] opacity-70">Refreshing every minute</span>
        </div>
      )}
      {loading ? (
        <div className="flex items-center justify-center py-24 rounded-xl border border-border/50 bg-card/50">
          <LoadingSpinner label="Fetching session results..." />
        </div>
      ) : error ? (
        <div className="rounded-xl border border-border/50 bg-card/50 p-12 text-center">
          <AlertCircle className="mx-auto mb-3 h-8 w-8 text-destructive/50" />
          <p className="text-muted-foreground">{error}</p>
          <button 
            onClick={onRetry}
            className="text-xs text-primary hover:underline mt-2"
          >
            Try Again
          </button>
        </div>
      ) : results && results.length > 0 ? (
        <ResultsTable
          results={results}
          sessionName={sessionData?.type || session}
        />
      ) : (
        <div className="rounded-xl border border-border/50 bg-card/50 p-12 text-center">
          <p className="text-muted-foreground">No results available for this session yet.</p>
        </div>
      )}
    </div>
  );
}
