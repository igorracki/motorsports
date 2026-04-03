import { ApiError } from "./api-error";
import { env } from "@/lib/env";

export class HttpClient {
  private static getBaseUrl() {
    if (typeof window === "undefined") {
      // Server-side (SSR): Use internal Docker network from validated env
      return `${env.BACKEND_URL}/api`;
    }
    // Client-side (Browser): Use relative path to be proxied by Next.js
    return "/api";
  }

  private readonly baseUrl: string;

  constructor() {
    this.baseUrl = HttpClient.getBaseUrl();
  }

  private getCookie(name: string): string | null {
    if (typeof document === "undefined") return null;
    const value = `; ${document.cookie}`;
    const parts = value.split(`; ${name}=`);
    if (parts.length === 2) return parts.pop()?.split(";").shift() || null;
    return null;
  }

  /**
   * Helper for fetch calls with error handling
   */
  async fetchJson<T>(path: string, options?: RequestInit): Promise<T> {
    const url = `${this.baseUrl}${path.startsWith("/") ? path : `/${path}`}`;
    try {
      const csrfToken = this.getCookie("csrf_token");
      const headers = new Headers(options?.headers);

      if (!headers.has("Content-Type")) {
        headers.set("Content-Type", "application/json");
      }

      // Add CSRF token for state-changing requests (POST, PUT, DELETE)
      if (
        csrfToken &&
        options?.method &&
        ["POST", "PUT", "DELETE", "PATCH"].includes(options.method.toUpperCase())
      ) {
        headers.set("X-CSRF-Token", csrfToken);
      }

      const response = await fetch(url, {
        ...options,
        credentials: options?.credentials || "include",
        headers,
        // Next.js specific cache configuration
        next: {
          revalidate: 60,
          ...((options as Record<string, unknown>)?.next || {}),
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
  }
}
