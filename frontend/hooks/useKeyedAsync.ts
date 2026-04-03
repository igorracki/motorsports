import { useState, useCallback, useRef, useEffect } from "react";
import { ErrorTranslator } from "@/lib/error-translator";

interface KeyedAsyncState<T> {
  data: Record<string, T>;
  loading: Record<string, boolean>;
  error: Record<string, string>;
  rawError: Record<string, unknown>;
}

/**
 * Keyed Async Decorator Hook
 * Extends the useAsync concept for dictionaries (e.g., fetching results per session code).
 */
export function useKeyedAsync<T, Args extends unknown[]>(
  asyncFn: (key: string, ...args: Args) => Promise<T>
) {
  const [state, setState] = useState<KeyedAsyncState<T>>({
    data: {},
    loading: {},
    error: {},
    rawError: {},
  });

  const asyncFnRef = useRef(asyncFn);

  useEffect(() => {
    asyncFnRef.current = asyncFn;
  }, [asyncFn]);

  const execute = useCallback(
    async (key: string, ...args: Args) => {
      // Check if already loading or loaded to prevent duplicate fetches
      setState((prev) => {
        if (prev.loading[key] || prev.data[key]) return prev;
        return {
          ...prev,
          loading: { ...prev.loading, [key]: true },
          error: { ...prev.error, [key]: "" },
          rawError: { ...prev.rawError, [key]: null },
        };
      });

      try {
        const data = await asyncFnRef.current(key, ...args);
        setState((prev) => ({
          ...prev,
          data: { ...prev.data, [key]: data },
          loading: { ...prev.loading, [key]: false },
        }));
        return data;
      } catch (err) {
        const message = ErrorTranslator.toDisplayMessage(err);
        setState((prev) => ({
          ...prev,
          loading: { ...prev.loading, [key]: false },
          error: { ...prev.error, [key]: message },
          rawError: { ...prev.rawError, [key]: err },
        }));
        throw err;
      }
    },
    []
  );

  return {
    ...state,
    execute,
    reset: (key: string) =>
      setState((prev) => {
        const newData = { ...prev.data };
        delete newData[key];
        const newLoading = { ...prev.loading };
        delete newLoading[key];
        const newError = { ...prev.error };
        delete newError[key];
        const newRawError = { ...prev.rawError };
        delete newRawError[key];
        return { data: newData, loading: newLoading, error: newError, rawError: newRawError };
      }),
  };
}
