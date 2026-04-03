import { useState, useCallback, useRef, useEffect } from "react";
import { ErrorTranslator } from "@/lib/error-translator";

interface AsyncState<T> {
  data: T | null;
  loading: boolean;
  error: string | null;
  rawError: unknown | null;
}

/**
 * Standardizes loading, error, and data states for async operations.
 * Eliminates repetitive try/catch/finally blocks in components.
 */
export function useAsync<T, Args extends unknown[]>(
  asyncFn: (...args: Args) => Promise<T>,
  options: {
    onSuccess?: (data: T) => void;
    onError?: (error: string) => void;
    immediate?: boolean;
  } = {}
) {
  const [state, setState] = useState<AsyncState<T>>({
    data: null,
    loading: false,
    error: null,
    rawError: null,
  });

  const asyncFnRef = useRef(asyncFn);
  const optionsRef = useRef(options);

  useEffect(() => {
    asyncFnRef.current = asyncFn;
    optionsRef.current = options;
  }, [asyncFn, options]);

  const execute = useCallback(
    async (...args: Args) => {
      setState((prev) => ({ ...prev, loading: true, error: null, rawError: null }));
      try {
        const data = await asyncFnRef.current(...args);
        setState({ data, loading: false, error: null, rawError: null });
        optionsRef.current.onSuccess?.(data);
        return data;
      } catch (err) {
        const message = ErrorTranslator.toDisplayMessage(err);
        setState({ data: null, loading: false, error: message, rawError: err });
        optionsRef.current.onError?.(message);
        throw err;
      }
    },
    []
  );

  return {
    ...state,
    execute,
    reset: () => setState({ data: null, loading: false, error: null, rawError: null }),
  };
}
