"use client";

import { useState } from "react";
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
import { Switch } from "@/components/ui/Switch";

interface ResultsTableProps {
  results: DriverResult[];
  sessionName: string;
}

export function ResultsTable({ results, sessionName }: ResultsTableProps) {
  const [showDetails, setShowDetails] = useState(false);
  const normalizedSession = sessionName.toLowerCase();
  const isQualifying = normalizedSession.includes("qualifying") || normalizedSession === "sq" || normalizedSession.includes("shootout");
  const isPractice = normalizedSession.includes("practice") || normalizedSession.includes("fp");
  const isRace = normalizedSession === "race" || normalizedSession === "sprint" || normalizedSession === "s" || (normalizedSession.includes("race") && !normalizedSession.includes("practice"));

  // Find the overall fastest lap in the session for highlighting
  const sessionBestLapMS = results.reduce((min, result) => {
    if (result.fastestLapMS && (min === null || result.fastestLapMS < min)) {
      return result.fastestLapMS;
    }
    return min;
  }, null as number | null);

  return (
    <div className="rounded-xl border border-border/50 bg-card overflow-hidden shadow-sm">
      <div className="border-b border-border/50 bg-secondary/50 px-4 py-3 flex items-center justify-between">
        <h3 className="font-semibold text-foreground">{sessionName} Results</h3>
        {isQualifying && (
          <div className="flex items-center gap-2 md:hidden">
            <span className="text-[10px] uppercase font-bold text-muted-foreground tracking-wider">
              Show Times
            </span>
            <Switch
              checked={showDetails}
              onCheckedChange={setShowDetails}
            />
          </div>
        )}
      </div>
      <div className="overflow-x-auto">
        <Table>
          <TableHeader>
            <TableRow className="hover:bg-transparent border-border/50">
              <TableHead className="w-16 text-center text-muted-foreground font-bold">Pos</TableHead>
              <TableHead className="text-muted-foreground font-bold">Driver</TableHead>
              <TableHead className="hidden sm:table-cell text-muted-foreground font-bold">Team</TableHead>
              {isQualifying && (
                <>
                  <TableHead className={cn(
                    "md:table-cell text-right text-muted-foreground font-bold",
                    showDetails ? "table-cell" : "hidden"
                  )}>Q1</TableHead>
                  <TableHead className={cn(
                    "md:table-cell text-right text-muted-foreground font-bold",
                    showDetails ? "table-cell" : "hidden"
                  )}>Q2</TableHead>
                  <TableHead className={cn(
                    "md:table-cell text-right text-muted-foreground font-bold",
                    showDetails ? "table-cell" : "hidden"
                  )}>Q3</TableHead>
                </>
              )}
              <TableHead className="text-right text-muted-foreground font-bold whitespace-nowrap">
                {isQualifying || isPractice || isRace ? "Best Lap" : "Time"}
              </TableHead>
              <TableHead className="text-right text-muted-foreground font-bold">
                {isRace ? "Total Time" : "Gap"}
              </TableHead>
              {isRace && (
                <TableHead className="w-16 text-right text-muted-foreground font-bold whitespace-nowrap">
                  +/-
                </TableHead>
              )}
            </TableRow>
          </TableHeader>
          <TableBody>
            {results.map((result, index) => {
              const displayPosition = result.position > 0 ? result.position : index + 1;
              
              let displayTime = result.totalTime || result.fastestLap || result.status;
              let displayGap = result.gap;

              const isFastestLap = sessionBestLapMS !== null && result.fastestLapMS === sessionBestLapMS;

              if (isPractice) {
                displayTime = result.fastestLap || result.status;
                displayGap = result.gap;
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
                  <TableCell className="font-medium text-foreground whitespace-nowrap">
                    {result.driver.fullName}
                  </TableCell>
                  <TableCell className="hidden text-muted-foreground sm:table-cell whitespace-nowrap">
                    {result.driver.teamName}
                  </TableCell>
                  {isQualifying && (
                    <>
                      <TableCell className={cn(
                        "md:table-cell text-right font-mono text-sm text-muted-foreground",
                        showDetails ? "table-cell" : "hidden"
                      )}>
                        {result.qualifying_details?.q1 || "-"}
                      </TableCell>
                      <TableCell className={cn(
                        "md:table-cell text-right font-mono text-sm text-muted-foreground",
                        showDetails ? "table-cell" : "hidden"
                      )}>
                        {result.qualifying_details?.q2 || "-"}
                      </TableCell>
                      <TableCell className={cn(
                        "md:table-cell text-right font-mono text-sm text-muted-foreground",
                        showDetails ? "table-cell" : "hidden"
                      )}>
                        {result.qualifying_details?.q3 || "-"}
                      </TableCell>
                    </>
                  )}
                  <TableCell 
                    className={cn(
                      "text-right font-mono text-sm",
                      isFastestLap ? "text-purple-400 font-bold" : "text-foreground"
                    )}
                  >
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
                  {isRace && (
                    <TableCell
                      className={cn(
                        "text-right font-mono text-sm",
                        (result.raceDetails?.positionsChange ?? 0) > 0 && "text-success",
                        (result.raceDetails?.positionsChange ?? 0) < 0 && "text-destructive",
                        (result.raceDetails?.positionsChange ?? 0) === 0 && "text-muted-foreground"
                      )}
                    >
                      {result.raceDetails && result.raceDetails.positionsChange !== 0
                        ? (result.raceDetails.positionsChange > 0 ? `+${result.raceDetails.positionsChange}` : result.raceDetails.positionsChange)
                        : "-"}
                    </TableCell>
                  )}
                </TableRow>
              );
            })}
          </TableBody>
        </Table>
      </div>
    </div>
  );
}
