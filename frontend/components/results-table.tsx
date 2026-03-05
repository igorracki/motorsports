"use client";

import type { DriverResult } from "@/types/f1";
import { cn } from "@/lib/utils";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/Table";

interface ResultsTableProps {
  results: DriverResult[];
  sessionName: string;
}

export function ResultsTable({ results, sessionName }: ResultsTableProps) {
  const isQualifying = sessionName.toLowerCase().includes("qualifying");
  const isPractice = sessionName.toLowerCase().includes("practice");
  const isRace = sessionName.toLowerCase() === "race";

  return (
    <div className="rounded-xl border border-border/50 bg-card overflow-hidden shadow-sm">
      <div className="border-b border-border/50 bg-secondary/50 px-4 py-3">
        <h3 className="font-semibold text-foreground">{sessionName} Results</h3>
      </div>
      <Table>
        <TableHeader>
          <TableRow className="hover:bg-transparent border-border/50">
            <TableHead className="w-16 text-center text-muted-foreground font-bold">Pos</TableHead>
            <TableHead className="text-muted-foreground font-bold">Driver</TableHead>
            <TableHead className="hidden sm:table-cell text-muted-foreground font-bold">Team</TableHead>
            <TableHead className="text-right text-muted-foreground font-bold">
              {isQualifying || isPractice || isRace ? "Best Lap" : "Time"}
            </TableHead>
            <TableHead className="text-right text-muted-foreground font-bold">
              {isPractice ? "Laps" : isRace ? "Total Time" : "Gap"}
            </TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {results.map((result, index) => {
            const displayPosition = result.position > 0 ? result.position : index + 1;
            
            let displayTime = result.totalTime || result.fastestLap || result.status;
            let displayGap = result.gap;

            if (isPractice) {
              displayTime = result.fastestLap || result.status;
              displayGap = result.laps.toString();
            } else if (isQualifying) {
              displayTime = result.fastestLap || result.status;
              displayGap = result.gap;
            } else if (isRace) {
              displayTime = result.fastestLap || result.status;
              displayGap = displayPosition === 1 ? (result.totalTime || "Finished") : result.gap;
            }

            return (
              <TableRow
                key={result.driver.id}
                className={cn(
                  "border-border/30",
                  index < 3 && displayPosition > 0 && "bg-primary/5"
                )}
              >
                <TableCell className="text-center">
                  <span
                    className={cn(
                      "inline-flex h-7 w-7 items-center justify-center rounded-full text-sm font-bold",
                      displayPosition === 1 && "bg-amber-500/20 text-amber-400",
                      displayPosition === 2 && "bg-slate-400/20 text-slate-300",
                      displayPosition === 3 && "bg-orange-600/20 text-orange-400",
                      displayPosition > 3 && "text-muted-foreground"
                    )}
                  >
                    {displayPosition}
                  </span>
                </TableCell>
                <TableCell className="font-medium text-foreground">
                  {result.driver.fullName}
                </TableCell>
                <TableCell className="hidden text-muted-foreground sm:table-cell">
                  {result.driver.teamName}
                </TableCell>
                <TableCell className="text-right font-mono text-sm text-foreground">
                  {displayTime}
                </TableCell>
                <TableCell
                  className={cn(
                    "text-right font-mono text-sm",
                    displayGap === "-" || (isRace && displayPosition === 1) ? "text-primary" : "text-muted-foreground"
                  )}
                >
                  {displayGap}
                </TableCell>
              </TableRow>
            );
          })}
        </TableBody>
      </Table>
    </div>
  );
}
