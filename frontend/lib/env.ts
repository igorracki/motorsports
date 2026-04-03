import { z } from "zod";

/**
 * Helper to treat empty strings as undefined, allowing Zod's .default() to trigger.
 */
const emptyToUndefined = (val: unknown) => (val === "" ? undefined : val);

/**
 * Environment variables schema for the application.
 */
const envSchema = z.object({
  /**
   * Backend API URL. Required for SSR.
   */
  BACKEND_URL: z.preprocess(
    emptyToUndefined,
    z.string().url().default("http://backend:8080")
  ),

  /**
   * Public Backend API URL. Used for client-side configuration.
   */
  NEXT_PUBLIC_BACKEND_URL: z.preprocess(
    emptyToUndefined,
    z.string().url().default("http://localhost:8080")
  ),

  /**
   * Environment mode.
   */
  NODE_ENV: z.enum(["development", "test", "production"]).default("development"),
});

/**
 * Validates and returns the application environment.
 * In a Next.js context, this runs both on server and client (for public vars).
 * Private variables like BACKEND_URL are only available on the server.
 */
const parseEnv = () => {
  try {
    return envSchema.parse({
      BACKEND_URL: process.env.BACKEND_URL,
      NEXT_PUBLIC_BACKEND_URL: process.env.NEXT_PUBLIC_BACKEND_URL,
      NODE_ENV: process.env.NODE_ENV,
    });
  } catch (error) {
    if (error instanceof z.ZodError) {
      const missingVars = error.issues.map((i) => i.path.join(".")).join(", ");
      throw new Error(`Invalid or missing environment variables: ${missingVars}`);
    }
    throw error;
  }
};

export const env = parseEnv();
