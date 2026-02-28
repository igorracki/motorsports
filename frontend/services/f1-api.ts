import { 
  RaceWeekend, 
  RaceWeekendSchema, 
  DriverInfo,
  DriverInfoSchema
} from "@/types/f1";
import * as dummyData from "@/lib/events-data";
import driversData from "@/data/drivers.json";
import { z } from "zod";

const _BASE_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080/api";

/**
 * Custom Error class for API related errors
 */
export class ApiError extends Error {
  constructor(
    public message: string,
    public status?: number,
    public code?: string
  ) {
    super(message);
    this.name = "ApiError";
  }
}

/**
 * Service to handle all F1 data fetching.
 */
export const f1Api = {
  /**
   * Helper for fetch calls with error handling
   */
  async fetchJson<T>(url: string, options?: RequestInit): Promise<T> {
    try {
      const response = await fetch(url, {
        ...options,
        headers: {
          "Content-Type": "application/json",
          ...options?.headers,
        },
      });

      if (!response.ok) {
        const errorData = await response.json().catch(() => ({}));
        throw new ApiError(
          errorData.message || `API error: ${response.statusText}`,
          response.status,
          errorData.error
        );
      }

      return await response.json();
    } catch (error) {
      if (error instanceof ApiError) throw error;
      throw new ApiError(
        error instanceof Error ? error.message : "Network error occurred"
      );
    }
  },

  /**
   * Fetches the full list of drivers
   */
  async getDrivers(): Promise<DriverInfo[]> {
    // Mapping the JSON data to our Schema
    const drivers = driversData.map(d => ({
      id: d.id,
      number: "0", // Dummy
      fullName: d.name,
      countryCode: "", // Dummy
      teamName: d.team
    }));

    return z.array(DriverInfoSchema).parse(drivers);
  },

  /**
   * Fetches the full schedule for a given year.
   */
  async getSchedule(year: number): Promise<RaceWeekend[]> {
    // Currently using dummy data for development
    const rawWeekends = dummyData.getEventsByYear(year);
    
    // Validate dummy data against our schema to ensure consistency
    const schedule = rawWeekends.map(raw => ({
      round: rawWeekends.indexOf(raw) + 1,
      fullName: raw.title,
      name: raw.title,
      location: raw.location,
      country: raw.country,
      countryCode: raw.countryCode,
      startDateMS: new Date().getTime(),
      sessions: raw.sessions.map(s => ({
        type: s.code,
        timeLocalMS: 0,
        timeUTCMS: 0,
        results: s.results?.map(r => ({
          position: r.position,
          driver: {
            id: r.driver.toLowerCase().replace(" ", "-"),
            number: "0",
            fullName: r.driver,
            countryCode: "",
            teamName: r.team
          },
          laps: 0,
          status: "Finished",
          gap: r.gap
        }))
      }))
    }));

    return z.array(RaceWeekendSchema).parse(schedule);
  },

  /**
   * Fetches details for a specific race weekend.
   */
  async getRaceWeekend(year: number, round: string): Promise<RaceWeekend | null> {
    const rawWeekends = dummyData.getEventsByYear(year);
    const roundIdx = parseInt(round, 10) - 1;
    const raw = rawWeekends[roundIdx];
    
    if (!raw) return null;

    const mapped: RaceWeekend = {
      round: roundIdx + 1,
      fullName: raw.title,
      name: raw.title,
      location: raw.location,
      country: raw.country,
      countryCode: raw.countryCode,
      startDateMS: 0,
      sessions: raw.sessions.map(s => ({
        type: s.code,
        timeLocalMS: 0,
        timeUTCMS: 0,
        results: s.results?.map(r => ({
          position: r.position,
          driver: {
            id: r.driver.toLowerCase().replace(" ", "-"),
            number: "0",
            fullName: r.driver,
            countryCode: "",
            teamName: r.team
          },
          laps: 0,
          status: "Finished",
          gap: r.gap
        }))
      }))
    };

    return RaceWeekendSchema.parse(mapped);
  }
};
