import { Metadata } from "next";
import { f1Api } from "@/services/f1-api";
import { RaceWeekendDashboard } from "@/components/features/RaceWeekendDashboard";
import { notFound } from "next/navigation";

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
    if (!raceWeekend) return { title: "Race Not Found" };

    return {
      title: `${raceWeekend.name} ${year}`,
      description: `Results, track details and predictions for the ${raceWeekend.fullName} at ${raceWeekend.location}.`,
    };
  } catch (error) {
    return { title: "Error" };
  }
}

export default async function RaceWeekendPage({ params }: RaceWeekendPageProps) {
  const { round, year: yearStr } = await params;
  const year = Number(yearStr) || 2026;
  
  const raceWeekend = await f1Api.getRaceWeekend(year, round);

  if (!raceWeekend) {
    notFound();
  }

  return <RaceWeekendDashboard raceWeekend={raceWeekend} year={year} />;
}
