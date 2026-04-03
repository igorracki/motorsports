"use client";

import React, { createContext, useContext, useMemo } from "react";
import { HttpClient } from "@/services/api/http-client";
import { createApiClients } from "@/services/api/factory";

type ApiContextValue = ReturnType<typeof createApiClients>;

const ApiContext = createContext<ApiContextValue | null>(null);

export function ApiProvider({ children }: { children: React.ReactNode }) {
  const api = useMemo(() => {
    const httpClient = new HttpClient();
    return createApiClients(httpClient);
  }, []);

  return <ApiContext.Provider value={api}>{children}</ApiContext.Provider>;
}

export function useApi() {
  const context = useContext(ApiContext);
  if (!context) {
    throw new Error("useApi must be used within an ApiProvider");
  }
  return context;
}
