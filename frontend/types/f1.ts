import { z } from "zod";

// --- Driver Schemas ---
export const DriverInfoSchema = z.object({
  id: z.string(),
  number: z.string(),
  fullName: z.string(),
  countryCode: z.string(),
  teamName: z.string(),
});

export type DriverInfo = z.infer<typeof DriverInfoSchema>;

// --- Session Result Schemas ---
export const RaceDetailsSchema = z.object({
  gridPosition: z.number(),
  status: z.string(),
  positionsChange: z.number(),
});

export type RaceDetails = z.infer<typeof RaceDetailsSchema>;

export const QualifyingDetailsSchema = z.object({
  q1MS: z.number().optional(),
  q1: z.string().optional(),
  q2MS: z.number().optional(),
  q2: z.string().optional(),
  q3MS: z.number().optional(),
  q3: z.string().optional(),
});

export type QualifyingDetails = z.infer<typeof QualifyingDetailsSchema>;

export const DriverResultSchema = z.object({
  position: z.number(),
  driver: DriverInfoSchema,
  laps: z.number(),
  status: z.string(),
  totalTimeMS: z.number().optional(),
  totalTime: z.string().optional(),
  gapMS: z.number().optional(),
  gap: z.string(),
  fastestLapMS: z.number().optional(),
  fastestLap: z.string().optional(),
  raceDetails: RaceDetailsSchema.optional(),
  qualifyingDetails: QualifyingDetailsSchema.optional(),
});

export type DriverResult = z.infer<typeof DriverResultSchema>;

// --- Race Weekend Schemas ---
export const SessionSchema = z.object({
  type: z.string(),
  timeLocalMS: z.number(),
  timeLocal: z.string().optional(),
  timeUTCMS: z.number(),
  timeUTC: z.string().optional(),
  results: z.array(DriverResultSchema).optional(),
});

export type Session = z.infer<typeof SessionSchema>;

export const RaceWeekendSchema = z.object({
  round: z.number(),
  fullName: z.string(),
  name: z.string(),
  location: z.string(),
  country: z.string(),
  countryCode: z.string().optional(), // Added for flag logic
  startDateMS: z.number(),
  startDate: z.string().optional(),
  sessions: z.array(SessionSchema),
});

export type RaceWeekend = z.infer<typeof RaceWeekendSchema>;

// --- Circuit Schemas ---
export const CircuitLayoutPointSchema = z.object({
  x: z.number(),
  y: z.number(),
});

export type CircuitLayoutPoint = z.infer<typeof CircuitLayoutPointSchema>;

export const CircuitSchema = z.object({
  circuitName: z.string(),
  location: z.string(),
  country: z.string(),
  latitude: z.number().optional(),
  longitude: z.number().optional(),
  lengthKM: z.number().optional(),
  corners: z.number().optional(),
  layout: z.array(CircuitLayoutPointSchema).optional(),
  eventName: z.string(),
  eventDateMS: z.number(),
  eventDate: z.string().optional(),
  rotation: z.number().optional(),
});

export type Circuit = z.infer<typeof CircuitSchema>;

// --- Prediction Schemas ---
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

// --- API Response Schemas ---
export const ScheduleResponseSchema = z.object({
  schedule: z.array(RaceWeekendSchema),
});

export type ScheduleResponse = z.infer<typeof ScheduleResponseSchema>;
