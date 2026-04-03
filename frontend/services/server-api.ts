import { HttpClient } from "./api/http-client";
import { createApiClients } from "./api/factory";

/**
 * Server-side API factory.
 * Instantiates the API clients for Server Components.
 */
export function getServerApi() {
  const httpClient = new HttpClient();
  return createApiClients(httpClient);
}
