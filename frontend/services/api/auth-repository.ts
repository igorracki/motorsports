import {
  LoginRequest,
  RegisterRequest,
  AuthResponse,
  AuthResponseSchema
} from "@/types/auth";
import { HttpClient } from "./http-client";

export class AuthRepository {
  constructor(private client: HttpClient) { }

  async login(request: LoginRequest): Promise<AuthResponse> {
    const data = await this.client.fetchJson<unknown>(`/auth/login`, {
      method: "POST",
      body: JSON.stringify(request),
    });
    return AuthResponseSchema.parse(data);
  }

  async register(request: RegisterRequest): Promise<AuthResponse> {
    const data = await this.client.fetchJson<unknown>(`/auth/register`, {
      method: "POST",
      body: JSON.stringify(request),
    });
    return AuthResponseSchema.parse(data);
  }

  async logout(): Promise<void> {
    await this.client.fetchJson<void>(`/auth/logout`, {
      method: "POST",
    });
  }

  async getMe(): Promise<AuthResponse> {
    const data = await this.client.fetchJson<unknown>(`/auth/me`);
    return AuthResponseSchema.parse(data);
  }
}
