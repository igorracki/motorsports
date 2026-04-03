import { z } from "zod";
import {
  FriendRequest,
  FriendRequestSchema,
  LeaderboardEntry,
  LeaderboardEntrySchema
} from "@/types/f1";
import { HttpClient } from "./http-client";

export class FriendRepository {
  constructor(private client: HttpClient) { }

  async sendFriendRequest(identifier: string): Promise<void> {
    await this.client.fetchJson<void>(`/users/friends/request`, {
      method: "POST",
      body: JSON.stringify({ identifier }),
    });
  }

  async getPendingRequests(): Promise<FriendRequest[]> {
    const data = await this.client.fetchJson<unknown[]>(`/users/friends/requests`);
    return z.array(FriendRequestSchema).parse(data || []);
  }

  async handleFriendRequest(requestId: string, action: "accept" | "deny"): Promise<void> {
    await this.client.fetchJson<void>(`/users/friends/requests/${requestId}`, {
      method: "PUT",
      body: JSON.stringify({ action }),
    });
  }

  async getLeaderboard(season: number): Promise<LeaderboardEntry[]> {
    const data = await this.client.fetchJson<unknown[]>(`/users/friends/leaderboard/${season}`);
    return z.array(LeaderboardEntrySchema).parse(data || []);
  }
}
