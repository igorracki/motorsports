import { HttpClient } from "./http-client";
import { AuthRepository } from "./auth-repository";
import { RaceRepository } from "./race-repository";
import { PredictionRepository } from "./prediction-repository";
import { FriendRepository } from "./friend-repository";
import { UserRepository } from "./user-repository";
import { ConfigRepository } from "./config-repository";

/**
 * Ensures we do not duplicate dependency injection setup across client and server.
 */
export function createApiClients(httpClient: HttpClient) {
  return {
    authRepo: new AuthRepository(httpClient),
    raceRepo: new RaceRepository(httpClient),
    predictionRepo: new PredictionRepository(httpClient),
    friendRepo: new FriendRepository(httpClient),
    userRepo: new UserRepository(httpClient),
    configRepo: new ConfigRepository(httpClient),
  };
}
