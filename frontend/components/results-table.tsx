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
  return (
    <div className="rounded-xl border border-border/50 bg-card overflow-hidden">
      <div className="border-b border-border/50 bg-secondary/50 px-4 py-3">
        <h3 className="font-semibold text-foreground">{sessionName} Results</h3>
      </div>
      <Table>
        <TableHeader>
          <TableRow className="hover:bg-transparent">
            <TableHead className="w-16 text-center">Pos</TableHead>
            <TableHead>Driver</TableHead>
            <TableHead className="hidden sm:table-cell">Team</TableHead>
            <TableHead className="text-right">Time</TableHead>
            <TableHead className="text-right">Gap</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {results.map((result, index) => (
            <TableRow
              key={result.position}
              className={cn(
                index < 3 && "bg-primary/5"
              )}
            >
              <TableCell className="text-center">
                <span
                  className={cn(
                    "inline-flex h-7 w-7 items-center justify-center rounded-full text-sm font-bold",
                    result.position === 1 && "bg-amber-500/20 text-amber-400",
                    result.position === 2 && "bg-slate-400/20 text-slate-300",
                    result.position === 3 && "bg-orange-600/20 text-orange-400",
                    result.position > 3 && "text-muted-foreground"
                  )}
                >
                  {result.position}
                </span>
              </TableCell>
              <TableCell className="font-medium text-foreground">
                {result.driver.fullName}
              </TableCell>
              <TableCell className="hidden text-muted-foreground sm:table-cell">
                {result.driver.teamName}
              </TableCell>
              <TableCell className="text-right font-mono text-sm text-foreground">
                {result.totalTime || result.status}
              </TableCell>
              <TableCell
                className={cn(
                  "text-right font-mono text-sm",
                  result.gap === "-" ? "text-primary" : "text-muted-foreground"
                )}
              >
                {result.gap}
              </TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>
    </div>
  );
}
