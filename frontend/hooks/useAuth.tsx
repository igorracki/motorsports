"use client";

import { createContext, useContext, useEffect, useState, ReactNode } from "react";
import { User, Profile } from "@/types/f1";
import { LoginRequest, RegisterRequest } from "@/types/auth";
import { useApi } from "@/components/providers/api-provider";
import { ErrorTranslator } from "@/lib/error-translator";

interface AuthContextType {
  user: User | null;
  profile: Profile | null;
  isAuthenticated: boolean;
  isLoading: boolean;
  login: (request: LoginRequest) => Promise<void>;
  register: (request: RegisterRequest) => Promise<void>;
  logout: () => Promise<void>;
  error: string | null;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export function AuthProvider({ children }: { children: ReactNode }) {
  const { authRepo } = useApi();
  const [user, setUser] = useState<User | null>(null);
  const [profile, setProfile] = useState<Profile | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    async function checkSession() {
      try {
        const response = await authRepo.getMe();
        
        if (response.user && response.user.id) {
            setUser(response.user);
            setProfile(response.profile);
        }
      } catch (err) {
        if (!ErrorTranslator.isSilent(err)) {
          console.debug("Session check failed", err);
        }
      } finally {
        setIsLoading(false);
      }
    }

    checkSession();
  }, [authRepo]);

  const login = async (request: LoginRequest) => {
    setError(null);
    try {
      const response = await authRepo.login(request);
      setUser(response.user);
      setProfile(response.profile);
    } catch (err) {
      const message = ErrorTranslator.toDisplayMessage(err);
      setError(message);
      throw err;
    }
  };

  const register = async (request: RegisterRequest) => {
    setError(null);
    try {
      const response = await authRepo.register(request);
      setUser(response.user);
      setProfile(response.profile);
    } catch (err) {
      const message = ErrorTranslator.toDisplayMessage(err);
      setError(message);
      throw err;
    }
  };

  const logout = async () => {
    try {
      await authRepo.logout();
    } catch (err) {
      console.error("Logout failed", err);
    } finally {
      setUser(null);
      setProfile(null);
    }
  };

  return (
    <AuthContext.Provider
      value={{
        user,
        profile,
        isAuthenticated: !!user,
        isLoading,
        login,
        register,
        logout,
        error,
      }}
    >
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error("useAuth must be used within an AuthProvider");
  }
  return context;
}
