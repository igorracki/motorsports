"use client";

import React, { createContext, useContext, useState, useEffect } from "react";
import { AppConfig } from "@/types/f1";
import { useApi } from "./api-provider";
import { PredictionPolicy } from "@/lib/policies/prediction-policy";

interface ConfigContextValue {
  config: AppConfig | null;
  loading: boolean;
  error: Error | null;
}

const ConfigContext = createContext<ConfigContextValue | null>(null);

export function ConfigProvider({ children }: { children: React.ReactNode }) {
  const { configRepo, predictionRepo } = useApi();
  const [config, setConfig] = useState<AppConfig | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<Error | null>(null);

  useEffect(() => {
    async function init() {
      try {
        const [appConfig, policyConfig] = await Promise.all([
          configRepo.fetchConfig(),
          predictionRepo.getPredictionPolicy()
        ]);
        
        setConfig(appConfig);
        PredictionPolicy.setConfiguration(policyConfig);
      } catch (err) {
        console.error("Failed to initialize application config:", err);
        setError(err instanceof Error ? err : new Error("Unknown error during initialization"));
      } finally {
        setLoading(false);
      }
    }

    init();
  }, [configRepo, predictionRepo]);

  return (
    <ConfigContext.Provider value={{ config, loading, error }}>
      {children}
    </ConfigContext.Provider>
  );
}

export function useConfig() {
  const context = useContext(ConfigContext);
  if (!context) {
    throw new Error("useConfig must be used within a ConfigProvider");
  }
  return context;
}
