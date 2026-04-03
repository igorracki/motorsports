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
    isPredicted: false, // UI tracking only
    correct: false,      // UI tracking only
    points: 0,          // UI tracking only
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
    gap: z.string().optional().nullable(),
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
    gap: data.gap ?? undefined,
    fastestLapMS: data.fastest_lap_ms ?? undefined,
    fastestLap: data.fastest_lap ?? undefined,
    raceDetails: data.race_details ?? undefined,
    qualifyingDetails: data.qualifying_details ?? undefined,
  }));

export type DriverResult = z.infer<typeof DriverResultSchema>;

export const SessionSchema = z
  .object({
    type: z.string(),
    session_code: z.string().optional(),
    time_local: z.string().optional().nullable(),
    time_utc_ms: z.number(),
    time_utc: z.string().optional().nullable(),
    utc_offset_ms: z.number(),
    is_locked: z.boolean().default(false),
    is_live: z.boolean().default(false),
    is_completed: z.boolean().default(false),
    results: z.array(DriverResultSchema).optional(),
  })
  .transform((data) => ({
    type: data.type,
    sessionCode: data.session_code,
    timeLocal: data.time_local ?? undefined,
    timeLocalMS: 0, // Fallback for now as backend doesn't provide it
    timeUTCMS: data.time_utc_ms,
    timeUTC: data.time_utc ?? undefined,
    utcOffsetMS: data.utc_offset_ms,
    isLocked: data.is_locked,
    isLive: data.is_live,
    isCompleted: data.is_completed,
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
    start_date_utc: z.string().optional().nullable(),
    start_date_utc_ms: z.number().optional(),
    end_date_local: z.string().optional().nullable(),
    end_date_local_ms: z.number().optional(),
    end_date_utc: z.string().optional().nullable(),
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
    startDateLocal: data.start_date_local ?? undefined,
    startDateUTC: data.start_date_utc ?? undefined,
    startDateUTCMS: data.start_date_utc_ms ?? 0,
    endDateLocal: data.end_date_local ?? undefined,
    endDateLocalMS: data.end_date_local_ms ?? 0,
    endDateUTC: data.end_date_utc ?? undefined,
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
    rotation: z.number(),
    max_speed_kmh: z.number(),
    max_altitude_m: z.number(),
    min_altitude_m: z.number(),
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

export const PredictionEntrySchema = z
  .object({
    prediction_id: z.string(),
    position: z.number(),
    driver_id: z.string(),
    correct: z.boolean().optional(),
  })
  .transform((data) => ({
    predictionId: data.prediction_id,
    position: data.position,
    driverId: data.driver_id,
    correct: data.correct,
  }));

export type PredictionEntry = z.infer<typeof PredictionEntrySchema>;

export const PredictionSchema = z
  .object({
    id: z.string(),
    user_id: z.string(),
    year: z.number(),
    round: z.number(),
    session_type: z.string(),
    score: z.number().optional().nullable(),
    created_at: z.string(),
    updated_at: z.string(),
    entries: z.array(PredictionEntrySchema),
  })
  .transform((data) => ({
    id: data.id,
    userId: data.user_id,
    year: data.year,
    round: data.round,
    sessionType: data.session_type,
    score: data.score ?? undefined,
    createdAt: data.created_at,
    updatedAt: data.updated_at,
    entries: data.entries,
  }));

export type Prediction = z.infer<typeof PredictionSchema>;

export const UserSchema = z.object({
  id: z.string(),
  email: z.string(),
  created_at: z.string(),
}).transform(data => ({
  id: data.id,
  email: data.email,
  createdAt: data.created_at,
}));

export type User = z.infer<typeof UserSchema>;

export const ProfileSchema = z.object({
  user_id: z.string(),
  display_name: z.string(),
}).transform(data => ({
  userId: data.user_id,
  displayName: data.display_name,
}));

export type Profile = z.infer<typeof ProfileSchema>;

export const UserScoreSchema = z.object({
  user_id: z.string(),
  score_type: z.string(),
  season: z.number().optional(),
  value: z.number(),
  updated_at: z.string().optional(),
}).transform(data => ({
  userId: data.user_id,
  scoreType: data.score_type,
  season: data.season,
  value: data.value,
  updatedAt: data.updated_at,
}));

export type UserScore = z.infer<typeof UserScoreSchema>;

export const UserProfileResponseSchema = z.object({
  user: UserSchema,
  profile: ProfileSchema,
  scores: z.array(UserScoreSchema),
});

export type UserProfileResponse = z.infer<typeof UserProfileResponseSchema>;

export const FriendRequestSchema = z.object({
  id: z.string(),
  sender_id: z.string(),
  receiver_id: z.string(),
  status: z.string(),
  created_at: z.string(),
  sender_email: z.string().optional(),
  sender_name: z.string().optional(),
}).transform(data => ({
  id: data.id,
  senderId: data.sender_id,
  receiverId: data.receiver_id,
  status: data.status,
  createdAt: data.created_at,
  senderEmail: data.sender_email,
  senderName: data.sender_name,
}));

export type FriendRequest = z.infer<typeof FriendRequestSchema>;

export const LeaderboardEntrySchema = z.object({
  position: z.number(),
  user_id: z.string(),
  display_name: z.string(),
  score: z.number(),
}).transform(data => ({
  position: data.position,
  userId: data.user_id,
  displayName: data.display_name,
  score: data.score,
}));

export type LeaderboardEntry = z.infer<typeof LeaderboardEntrySchema>;

export const ScheduleResponseSchema = z.object({
  schedule: z.array(RaceWeekendSchema),
});

export type ScheduleResponse = z.infer<typeof ScheduleResponseSchema>;

export const PositionPointsSchema = z.object({
  position: z.number(),
  points: z.number(),
});

export type PositionPoints = z.infer<typeof PositionPointsSchema>;

export const SessionScoringRulesSchema = z.object({
  session_type: z.string(),
  rules: z.array(PositionPointsSchema),
}).transform(data => ({
  sessionType: data.session_type,
  rules: data.rules,
}));

export type SessionScoringRules = z.infer<typeof SessionScoringRulesSchema>;

export const PredictionPolicyConfigSchema = z.object({
  lock_threshold_ms: z.number(),
  pre_session_buffer_ms: z.number(),
  session_duration_ms: z.number(),
  revalidation_window_ms: z.number(),
}).transform(data => ({
  lockThresholdMS: data.lock_threshold_ms,
  preSessionBufferMS: data.pre_session_buffer_ms,
  sessionDurationMS: data.session_duration_ms,
  revalidationWindowMS: data.revalidation_window_ms,
}));

export type PredictionPolicyConfig = z.infer<typeof PredictionPolicyConfigSchema>;

export const DriverMetadataSchema = z.object({
  id: z.string(),
  full_name: z.string(),
  team_name: z.string(),
  team_color: z.string(),
  country_code: z.string(),
}).transform(data => ({
  id: data.id,
  fullName: data.full_name,
  teamName: data.team_name,
  teamColor: data.team_color,
  countryCode: data.country_code,
}));

export type DriverMetadata = z.infer<typeof DriverMetadataSchema>;

export const ValidationConfigSchema = z.object({
  min_year: z.number(),
  max_year: z.number(),
  min_round: z.number(),
  max_round: z.number(),
  min_entries: z.number(),
  max_entries: z.number(),
}).transform(data => ({
  minYear: data.min_year,
  maxYear: data.max_year,
  minRound: data.min_round,
  maxRound: data.max_round,
  minEntries: data.min_entries,
  maxEntries: data.max_entries,
}));

export type ValidationConfig = z.infer<typeof ValidationConfigSchema>;

export const AppConfigSchema = z.object({
  drivers: z.array(DriverMetadataSchema),
  session_mappings: z.record(z.string()),
  validation: ValidationConfigSchema,
}).transform(data => ({
  drivers: data.drivers,
  sessionMappings: data.session_mappings,
  validation: data.validation,
}));

export type AppConfig = z.infer<typeof AppConfigSchema>;

export interface SubmitPredictionRequest {
  year: number;
  round: number;
  session_type: string;
  entries: {
    position: number;
    driver_id: string;
  }[];
}
