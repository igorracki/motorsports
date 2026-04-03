import { z } from "zod";
import {
  Prediction,
  PredictionSchema,
  SubmitPredictionRequest,
  SessionScoringRules,
  SessionScoringRulesSchema
} from "@/types/f1";
import { HttpClient } from "./http-client";

export class PredictionRepository {
  constructor(private client: HttpClient) { }

  async getRoundPredictions(userId: string, year: number, round: number): Promise<Prediction[]> {
    const data = await this.client.fetchJson<unknown[]>(
      `/users/${userId}/predictions/${year}/${round}`
    );
    return z.array(PredictionSchema).parse(data || []);
  }

  async submitPrediction(userId: string, prediction: SubmitPredictionRequest): Promise<Prediction> {
    const data = await this.client.fetchJson<unknown>(`/users/${userId}/predictions`, {
      method: "POST",
      body: JSON.stringify(prediction),
    });
    return PredictionSchema.parse(data);
  }

  async getScoringRules(): Promise<SessionScoringRules[]> {
    const data = await this.client.fetchJson<unknown[]>(`/predictions/scoring-rules`, {
      next: { revalidate: 86400 } // Revalidate every 24 hours
    });
    return z.array(SessionScoringRulesSchema).parse(data || []);
  }
}
