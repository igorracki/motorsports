import { Metadata } from "next";
import { getServerApi } from "@/services/server-api";
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
  
  const { raceRepo } = getServerApi();
  try {
    const raceWeekend = await raceRepo.getRaceWeekend(year, round);
    if (!raceWeekend) return { title: "Race Not Found" };

    return {
      title: `${raceWeekend.name} ${year}`,
      description: `Results, track details and predictions for the ${raceWeekend.fullName} at ${raceWeekend.location}.`,
    };
  } catch {
    return { title: "Error" };
  }
}

export default async function RaceWeekendPage({ params }: RaceWeekendPageProps) {
  const { round, year: yearStr } = await params;
  const year = Number(yearStr) || 2026;
  
  const { raceRepo } = getServerApi();
  const raceWeekend = await raceRepo.getRaceWeekend(year, round);

  if (!raceWeekend) {
    notFound();
  }

  return (
    <RaceWeekendDashboard 
      raceWeekend={raceWeekend} 
      year={year} 
      serverTime={new Date().getTime()}
    />
  );
}
