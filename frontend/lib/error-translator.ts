import { ApiError } from "@/services/api/api-error";

/**
 * Centralizes the translation of Domain Errors into user-friendly messages.
 */
export const ErrorTranslator = {
  /**
   * Translates an ApiError or generic Error into a displayable string.
   */
  toDisplayMessage(error: unknown): string {
    if (!(error instanceof ApiError)) {
      if (error instanceof Error) return error.message;
      return "An unexpected error occurred. Please try again.";
    }

    switch (error.code) {
      case "UNAUTHORIZED":
        return "Your session has expired. Please log in again.";
      case "FORBIDDEN":
        return "You do not have permission to perform this action.";
      case "NOT_FOUND":
        return "The requested information could not be found.";
      case "VALIDATION_ERROR":
        return "Please check your input and try again.";
      case "NETWORK_ERROR":
        return "Unable to connect to the server. Please check your internet connection.";
      default:
        if (error.status) {
          if (error.status === 401) return "Please log in to continue.";
          if (error.status === 403) return "Access denied.";
          if (error.status === 404) return "Resource not found.";
          if (error.status >= 500) return "Server error. Our team has been notified.";
        }
        return error.message || "Something went wrong.";
    }
  },

  /**
   * Determines if an error is a "soft" error that can be ignored or handled silently.
   */
  isSilent(error: unknown): boolean {
    if (error instanceof ApiError && error.status === 401) {
      return true;
    }
    return false;
  }
};
