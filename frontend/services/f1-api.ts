import { 
  RaceWeekend, 
  RaceWeekendSchema, 
  DriverInfo,
  DriverInfoSchema,
  Circuit,
  CircuitSchema,
  DriverResult,
  DriverResultSchema,
  UserProfileResponse,
  UserProfileResponseSchema,
  Prediction,
  PredictionSchema,
  SessionScoringRules,
  SessionScoringRulesSchema,
  FriendRequest,
  FriendRequestSchema,
  LeaderboardEntry,
  LeaderboardEntrySchema
} from "@/types/f1";
import { 
  LoginRequest, 
  RegisterRequest, 
  AuthResponse, 
  AuthResponseSchema 
} from "@/types/auth";
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
        credentials: options?.credentials || "include",
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

      // Handle empty responses (like 204 No Content or empty 200 OK)
      const text = await response.text();
      if (!text) {
        return {} as T;
      }

      return JSON.parse(text);
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
  },

  /**
   * Auth: Login
   */
  async login(request: LoginRequest): Promise<AuthResponse> {
    const data = await this.fetchJson<any>(`${_BASE_URL}/auth/login`, {
      method: "POST",
      body: JSON.stringify(request),
    });
    return AuthResponseSchema.parse(data);
  },

  /**
   * Auth: Register
   */
  async register(request: RegisterRequest): Promise<AuthResponse> {
    const data = await this.fetchJson<any>(`${_BASE_URL}/auth/register`, {
      method: "POST",
      body: JSON.stringify(request),
    });
    return AuthResponseSchema.parse(data);
  },

  /**
   * Auth: Logout
   */
  async logout(): Promise<void> {
    await this.fetchJson<void>(`${_BASE_URL}/auth/logout`, {
      method: "POST",
    });
  },

  /**
   * Auth: Get current user
   */
  async getMe(): Promise<AuthResponse> {
    const data = await this.fetchJson<any>(`${_BASE_URL}/auth/me`);
    return AuthResponseSchema.parse(data);
  },

  /**
   * User: Get full profile
   */
  async getUserProfile(userId: string): Promise<UserProfileResponse> {
    const data = await this.fetchJson<any>(`${_BASE_URL}/users/${userId}`);
    return UserProfileResponseSchema.parse(data);
  },

  /**
   * User: Get season scores/stats
   */
  async getSeasonScores(userId: string): Promise<any[]> {
    const data = await this.fetchJson<any[]>(`${_BASE_URL}/users/${userId}/stats/seasons`);
    return data || [];
  },

  /**
   * Predictions: Get predictions for a specific round
   */
  async getRoundPredictions(userId: string, year: number, round: number): Promise<Prediction[]> {
    const data = await this.fetchJson<any[]>(
      `${_BASE_URL}/users/${userId}/predictions/${year}/${round}`
    );
    return z.array(PredictionSchema).parse(data || []);
  },

  /**
   * Predictions: Submit a prediction
   */
  async submitPrediction(userId: string, prediction: Partial<Prediction>): Promise<Prediction> {
    const data = await this.fetchJson<any>(`${_BASE_URL}/users/${userId}/predictions`, {
      method: "POST",
      body: JSON.stringify(prediction),
    });
    return PredictionSchema.parse(data);
  },

  /**
   * Predictions: Get scoring rules
   */
  async getScoringRules(): Promise<SessionScoringRules[]> {
    const data = await this.fetchJson<any[]>(`${_BASE_URL}/predictions/scoring-rules`);
    return z.array(SessionScoringRulesSchema).parse(data || []);
  },

  /**
   * Friends: Send friend request
   */
  async sendFriendRequest(identifier: string): Promise<void> {
    await this.fetchJson<void>(`${_BASE_URL}/users/friends/request`, {
      method: "POST",
      body: JSON.stringify({ identifier }),
    });
  },

  /**
   * Friends: Get pending requests
   */
  async getPendingRequests(): Promise<FriendRequest[]> {
    const data = await this.fetchJson<any[]>(`${_BASE_URL}/users/friends/requests`);
    return z.array(FriendRequestSchema).parse(data || []);
  },

  /**
   * Friends: Handle friend request
   */
  async handleFriendRequest(requestId: string, action: "accept" | "deny"): Promise<void> {
    await this.fetchJson<void>(`${_BASE_URL}/users/friends/requests/${requestId}`, {
      method: "PUT",
      body: JSON.stringify({ action }),
    });
  },

  /**
   * Leaderboard: Get leaderboard for a season
   */
  async getLeaderboard(season: number): Promise<LeaderboardEntry[]> {
    const data = await this.fetchJson<any[]>(`${_BASE_URL}/users/friends/leaderboard/${season}`);
    return z.array(LeaderboardEntrySchema).parse(data || []);
  }
};
