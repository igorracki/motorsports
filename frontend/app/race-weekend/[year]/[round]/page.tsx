import { Metadata } from "next";
import { f1Api } from "@/services/f1-api";
import { RaceWeekendDashboard } from "@/components/features/RaceWeekendDashboard";
import Link from "next/link";

interface RaceWeekendPageProps {
  params: Promise<{ round: string; year: string }>;
}

export async function generateMetadata({
  params,
}: RaceWeekendPageProps): Promise<Metadata> {
  const { round, year: yearStr } = await params;
  const year = Number(yearStr) || 2026;
  
  try {
    const raceWeekend = await f1Api.getRaceWeekend(year, round);
    if (!raceWeekend) return { title: "Race Not Found | F1 Data Hub" };

    return {
      title: `${raceWeekend.name} ${year} | F1 Data Hub`,
      description: `Results, track details and predictions for the ${raceWeekend.fullName} at ${raceWeekend.location}.`,
    };
  } catch (error) {
    return { title: "Error | F1 Data Hub" };
  }
}

export default async function RaceWeekendPage({ params }: RaceWeekendPageProps) {
  const { round, year: yearStr } = await params;
  const year = Number(yearStr) || 2026;
  
  const raceWeekend = await f1Api.getRaceWeekend(year, round);

  if (!raceWeekend) {
    return (
      <div className="flex min-h-screen items-center justify-center bg-background">
        <div className="text-center">
          <h1 className="mb-4 text-2xl font-bold text-foreground">
            Race weekend not found
          </h1>
          <Link
            href="/"
            className="text-primary hover:text-primary/80 transition-colors"
          >
            Return to calendar
          </Link>
        </div>
      </div>
    );
  }

  return <RaceWeekendDashboard raceWeekend={raceWeekend} year={year} />;
}
