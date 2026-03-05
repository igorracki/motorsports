import { 
  RaceWeekend, 
  RaceWeekendSchema, 
  DriverInfo,
  DriverInfoSchema,
  Circuit,
  CircuitSchema,
  DriverResult,
  DriverResultSchema
} from "@/types/f1";
import { z } from "zod";

const getBaseUrl = () => {
  if (typeof window === "undefined") {
    // Server-side (SSR): Use internal Docker network
    const backendUrl = process.env.BACKEND_URL || "http://backend:8080";
    return `${backendUrl}/api`;
  }
  // Client-side (Browser): Use public-facing URL
  return process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080/api";
};

const _BASE_URL = getBaseUrl();

/**
 * Custom Error class for API related errors
 */
export class ApiError extends Error {
  constructor(
    public message: string,
    public status?: number,
    public code?: string,
    public url?: string
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
          errorData.error,
          url
        );
      }

      return await response.json();
    } catch (error) {
      if (error instanceof ApiError) throw error;
      
      // Log more details for network errors, especially on the server
      console.error(`Fetch failed for URL: ${url}`, error);
      
      throw new ApiError(
        error instanceof Error ? error.message : "Network error occurred",
        undefined,
        undefined,
        url
      );
    }
  },

  /**
   * Fetches the full list of drivers for a given year and round from the API.
   */
  async getDrivers(year: number, round: number): Promise<DriverInfo[]> {
    const data = await this.fetchJson<any[]>(`${_BASE_URL}/schedule/${year}/${round}/drivers`);
    return z.array(DriverInfoSchema).parse(data || []);
  },

  /**
   * Fetches the full schedule for a given year.
   */
  async getSchedule(year: number): Promise<RaceWeekend[]> {
    const data = await this.fetchJson<{ schedule: any[] }>(`${_BASE_URL}/schedule/${year}`);
    return z.array(RaceWeekendSchema).parse(data.schedule);
  },

  /**
   * Fetches details for a specific race weekend.
   */
  async getRaceWeekend(year: number, round: string): Promise<RaceWeekend | null> {
    const schedule = await this.getSchedule(year);
    const roundNumber = parseInt(round, 10);
    const raceWeekend = schedule.find(rw => rw.round === roundNumber);
    return raceWeekend || null;
  },

  /**
   * Fetches results for a specific session.
   */
  async getSessionResults(year: number, round: number, sessionCode: string): Promise<DriverResult[]> {
    const data = await this.fetchJson<any>(
      `${_BASE_URL}/schedule/${year}/${round}/${sessionCode}/results`
    );
    // The Backend returns a SessionResults object
    return z.array(DriverResultSchema).parse(data.results || []);
  },

  /**
   * Fetches circuit details.
   */
  async getCircuit(year: number, round: number): Promise<Circuit> {
    const data = await this.fetchJson<any>(`${_BASE_URL}/schedule/${year}/${round}/circuit`);
    return CircuitSchema.parse(data);
  }
};
