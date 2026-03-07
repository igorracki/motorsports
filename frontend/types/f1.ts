import { z } from "zod";

export const DriverInfoSchema = z
  .object({
    id: z.string(),
    number: z.string(),
    full_name: z.string(),
    country_code: z.string(),
    team_name: z.string(),
  })
  .transform((data) => ({
    id: data.id,
    number: data.number,
    fullName: data.full_name,
    countryCode: data.country_code,
    teamName: data.team_name,
  }));

export type DriverInfo = z.infer<typeof DriverInfoSchema>;

export const RaceDetailsSchema = z
  .object({
    grid_position: z.number(),
    status: z.string(),
    positions_change: z.number(),
  })
  .transform((data) => ({
    gridPosition: data.grid_position,
    status: data.status,
    positionsChange: data.positions_change,
  }));

export type RaceDetails = z.infer<typeof RaceDetailsSchema>;

export const QualifyingDetailsSchema = z
  .object({
    q1_ms: z.number().optional().nullable(),
    q1: z.string().optional().nullable(),
    q2_ms: z.number().optional().nullable(),
    q2: z.string().optional().nullable(),
    q3_ms: z.number().optional().nullable(),
    q3: z.string().optional().nullable(),
  })
  .transform((data) => ({
    q1MS: data.q1_ms ?? undefined,
    q1: data.q1 ?? undefined,
    q2MS: data.q2_ms ?? undefined,
    q2: data.q2 ?? undefined,
    q3MS: data.q3_ms ?? undefined,
    q3: data.q3 ?? undefined,
  }));

export type QualifyingDetails = z.infer<typeof QualifyingDetailsSchema>;

export const DriverResultSchema = z
  .object({
    position: z.number(),
    driver: DriverInfoSchema,
    laps: z.number(),
    status: z.string(),
    total_time_ms: z.number().optional().nullable(),
    total_time: z.string().optional().nullable(),
    gap_ms: z.number().optional().nullable(),
    gap: z.string(),
    fastest_lap_ms: z.number().optional().nullable(),
    fastest_lap: z.string().optional().nullable(),
    race_details: RaceDetailsSchema.optional().nullable(),
    qualifying_details: QualifyingDetailsSchema.optional().nullable(),
  })
  .transform((data) => ({
    position: data.position,
    driver: data.driver,
    laps: data.laps,
    status: data.status,
    totalTimeMS: data.total_time_ms ?? undefined,
    totalTime: data.total_time ?? undefined,
    gapMS: data.gap_ms ?? undefined,
    gap: data.gap,
    fastestLapMS: data.fastest_lap_ms ?? undefined,
    fastestLap: data.fastest_lap ?? undefined,
    raceDetails: data.race_details ?? undefined,
    qualifying_details: data.qualifying_details ?? undefined,
  }));

export type DriverResult = z.infer<typeof DriverResultSchema>;

export const SessionSchema = z
  .object({
    type: z.string(),
    session_code: z.string().optional(),
    time_local: z.string().optional().nullable(),
    time_utc_ms: z.number(),
    time_utc: z.string().optional().nullable(),
    results: z.array(DriverResultSchema).optional(),
  })
  .transform((data) => ({
    type: data.type,
    sessionCode: data.session_code,
    timeLocal: data.time_local ?? undefined,
    timeLocalMS: 0, // Fallback for now as backend doesn't provide it yet
    timeUTCMS: data.time_utc_ms,
    timeUTC: data.time_utc ?? undefined,
    results: data.results,
  }));

export type Session = z.infer<typeof SessionSchema>;

export const RaceWeekendSchema = z
  .object({
    round: z.number(),
    full_name: z.string(),
    name: z.string(),
    location: z.string(),
    country: z.string(),
    country_code: z.string().optional(),
    event_format: z.string().optional(),
    start_date_local_ms: z.number(),
    start_date_local: z.string().optional().nullable(),
    start_date_utc_ms: z.number().optional(),
    end_date_utc_ms: z.number().optional(),
    sessions: z.array(SessionSchema),
  })
  .transform((data) => ({
    round: data.round,
    fullName: data.full_name,
    name: data.name,
    location: data.location,
    country: data.country,
    countryCode: data.country_code,
    eventFormat: data.event_format,
    startDateMS: data.start_date_local_ms,
    startDate: data.start_date_local ?? undefined,
    startDateUTCMS: data.start_date_utc_ms ?? 0,
    endDateUTCMS: data.end_date_utc_ms ?? 0,
    sessions: data.sessions,
  }));

export type RaceWeekend = z.infer<typeof RaceWeekendSchema>;

export const CircuitLayoutPointSchema = z.object({
  x: z.number(),
  y: z.number(),
});

export type CircuitLayoutPoint = z.infer<typeof CircuitLayoutPointSchema>;

export const CircuitSchema = z
  .object({
    circuit_name: z.string(),
    location: z.string(),
    country: z.string(),
    latitude: z.number().optional().nullable(),
    longitude: z.number().optional().nullable(),
    length_km: z.number().optional().nullable(),
    corners: z.number().optional().nullable(),
    layout: z.array(CircuitLayoutPointSchema).optional(),
    event_name: z.string(),
    event_date_ms: z.number(),
    event_date: z.string().optional().nullable(),
    rotation: z.number().optional(),
    max_speed_kmh: z.number().optional(),
    max_altitude_m: z.number().optional(),
    min_altitude_m: z.number().optional(),
  })
  .transform((data) => ({
    circuitName: data.circuit_name,
    location: data.location,
    country: data.country,
    latitude: data.latitude ?? undefined,
    longitude: data.longitude ?? undefined,
    lengthKM: data.length_km ?? undefined,
    corners: data.corners ?? undefined,
    layout: data.layout,
    eventName: data.event_name,
    eventDateMS: data.event_date_ms,
    eventDate: data.event_date ?? undefined,
    rotation: data.rotation,
    maxSpeedKmh: data.max_speed_kmh,
    maxAltitudeM: data.max_altitude_m,
    minAltitudeM: data.min_altitude_m,
  }));

export type Circuit = z.infer<typeof CircuitSchema>;

export const PredictionEntrySchema = z.object({
  predictionId: z.string(),
  position: z.number(),
  driverId: z.string(),
});

export type PredictionEntry = z.infer<typeof PredictionEntrySchema>;

export const PredictionSchema = z.object({
  id: z.string(),
  userId: z.string(),
  year: z.number(),
  round: z.number(),
  sessionType: z.string(),
  score: z.number().optional(),
  createdAt: z.string(),
  updatedAt: z.string(),
  entries: z.array(PredictionEntrySchema),
});

export type Prediction = z.infer<typeof PredictionSchema>;

export const ScheduleResponseSchema = z.object({
  schedule: z.array(RaceWeekendSchema),
});

export type ScheduleResponse = z.infer<typeof ScheduleResponseSchema>;
