import { z } from "zod";
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
import { HttpClient } from "./http-client";

// Module-level singleton cache to persist across SSR requests
const globalScheduleCache = new Map<number, Promise<RaceWeekend[]>>();

export class RaceRepository {
  constructor(private client: HttpClient) { }

  async getDrivers(year: number, round: number): Promise<DriverInfo[]> {
    const data = await this.client.fetchJson<unknown>(`/schedule/${year}/${round}/drivers`, {
      next: { revalidate: 3600 } // Revalidate every hour
    });
    return z.array(DriverInfoSchema).parse(data || []);
  }

  async getSchedule(year: number): Promise<RaceWeekend[]> {
    if (globalScheduleCache.has(year)) {
      return globalScheduleCache.get(year)!;
    }

    const promise = this.client.fetchJson<{ schedule: unknown[] }>(`/schedule/${year}`, {
      next: { revalidate: 3600 } // Revalidate every hour
    }).then(data => z.array(RaceWeekendSchema).parse(data.schedule));

    globalScheduleCache.set(year, promise);

    // In case of error, remove from cache so it can be retried
    promise.catch(() => {
      if (globalScheduleCache.get(year) === promise) {
        globalScheduleCache.delete(year);
      }
    });

    return promise;
  }

  async getRaceWeekend(year: number, round: string): Promise<RaceWeekend | null> {
    const schedule = await this.getSchedule(year);
    const roundNumber = parseInt(round, 10);
    const raceWeekend = schedule.find(rw => rw.round === roundNumber);
    return raceWeekend || null;
  }

  async getSessionResults(year: number, round: number, sessionCode: string): Promise<DriverResult[]> {
    const data = await this.client.fetchJson<{ results: unknown[] }>(
      `/schedule/${year}/${round}/${sessionCode}/results`
    );
    return z.array(DriverResultSchema).parse(data.results || []);
  }

  async getCircuit(year: number, round: number): Promise<Circuit> {
    const data = await this.client.fetchJson<unknown>(`/schedule/${year}/${round}/circuit`, {
      next: { revalidate: 86400 } // Revalidate every 24 hours
    });
    return CircuitSchema.parse(data);
  }
}
