import Link from "next/link";
import Image from "next/image";
import { ArrowLeft } from "lucide-react";
import { StatusBadge } from "@/components/ui/StatusBadge";
import type { RaceWeekend } from "@/types/f1";
import type { RaceStatus } from "@/lib/race-utils";

interface DashboardHeaderProps {
  raceWeekend: RaceWeekend;
  year: number;
  status: RaceStatus;
}

export function DashboardHeader({ raceWeekend, year, status }: DashboardHeaderProps) {
  return (
    <header className="sticky top-14 z-50 border-b border-border/50 bg-background/95 backdrop-blur-sm">
      <div className="mx-auto flex max-w-7xl items-center gap-4 px-4 py-4 sm:px-6 lg:px-8">
        <Link
          href={`/calendar/${year}`}
          className="flex items-center gap-2 text-muted-foreground transition-colors hover:text-foreground"
        >
          <ArrowLeft className="h-5 w-5" />
          <span className="hidden sm:inline">Back to Calendar</span>
        </Link>
        <div className="h-6 w-px bg-border/50" />
        <div className="flex items-center gap-3">
          {raceWeekend.countryCode && (
            <Image
              src={`https://flagcdn.com/w80/${raceWeekend.countryCode.toLowerCase()}.png`}
              alt={`${raceWeekend.country} flag`}
              width={32}
              height={20}
              className="rounded-sm object-cover shadow-sm"
              unoptimized
            />
          )}
          <h1 className="text-lg font-bold sm:text-xl">
            {raceWeekend.fullName}
          </h1>
        </div>
        <StatusBadge status={status} className="ml-auto" />
      </div>
    </header>
  );
}
