"use client";

import type { SessionResult } from "@/lib/events-data";
import { cn } from "@/lib/utils";

interface ResultsTableProps {
  results: SessionResult[];
  sessionName: string;
}

export function ResultsTable({ results, sessionName }: ResultsTableProps) {
  return (
    <div className="overflow-hidden rounded-xl border border-border/50 bg-card">
      <div className="border-b border-border/50 bg-secondary/50 px-4 py-3">
        <h3 className="font-semibold text-foreground">{sessionName} Results</h3>
      </div>
      <div className="overflow-x-auto">
        <table className="w-full">
          <thead>
            <tr className="border-b border-border/50 text-left text-xs font-medium uppercase tracking-wider text-muted-foreground">
              <th className="px-4 py-3 text-center">Pos</th>
              <th className="px-4 py-3">Driver</th>
              <th className="hidden px-4 py-3 sm:table-cell">Team</th>
              <th className="px-4 py-3 text-right">Time</th>
              <th className="px-4 py-3 text-right">Gap</th>
            </tr>
          </thead>
          <tbody>
            {results.map((result, index) => (
              <tr
                key={result.position}
                className={cn(
                  "border-b border-border/30 transition-colors hover:bg-secondary/30",
                  index < 3 && "bg-primary/5"
                )}
              >
                <td className="px-4 py-3 text-center">
                  <span
                    className={cn(
                      "inline-flex h-7 w-7 items-center justify-center rounded-full text-sm font-bold",
                      result.position === 1 &&
                        "bg-amber-500/20 text-amber-400",
                      result.position === 2 &&
                        "bg-slate-400/20 text-slate-300",
                      result.position === 3 &&
                        "bg-orange-600/20 text-orange-400",
                      result.position > 3 && "text-muted-foreground"
                    )}
                  >
                    {result.position}
                  </span>
                </td>
                <td className="px-4 py-3 font-medium text-foreground">
                  {result.driver}
                </td>
                <td className="hidden px-4 py-3 text-muted-foreground sm:table-cell">
                  {result.team}
                </td>
                <td className="px-4 py-3 text-right font-mono text-sm text-foreground">
                  {result.time}
                </td>
                <td
                  className={cn(
                    "px-4 py-3 text-right font-mono text-sm",
                    result.gap === "-"
                      ? "text-primary"
                      : "text-muted-foreground"
                  )}
                >
                  {result.gap}
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
}
