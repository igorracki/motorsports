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
