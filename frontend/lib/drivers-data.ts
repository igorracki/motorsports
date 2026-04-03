import { DriverMetadata } from "@/types/f1";

export interface Driver {
  id: string;
  name: string;
  team: string;
  teamColor: string;
}

export function mapDriverMetadata(metadata: DriverMetadata): Driver {
  return {
    id: metadata.id.toLowerCase(),
    name: metadata.fullName,
    team: metadata.teamName,
    teamColor: metadata.teamColor,
  };
}
