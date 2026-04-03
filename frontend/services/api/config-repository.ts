import { AppConfig, AppConfigSchema } from "@/types/f1";
import { HttpClient } from "./http-client";

export class ConfigRepository {
  constructor(private client: HttpClient) { }

  async fetchConfig(): Promise<AppConfig> {
    const data = await this.client.fetchJson<unknown>(`/config`, {
      next: { revalidate: 86400 } // Revalidate every 24 hours
    });
    return AppConfigSchema.parse(data);
  }
}
