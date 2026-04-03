import { UserProfileResponse, UserProfileResponseSchema } from "@/types/f1";
import { HttpClient } from "./http-client";

export class UserRepository {
  constructor(private client: HttpClient) { }

  async getUserProfile(userId: string): Promise<UserProfileResponse> {
    const data = await this.client.fetchJson<unknown>(`/users/${userId}`);
    return UserProfileResponseSchema.parse(data);
  }
}
