"use client";

import React, { createContext, useContext, useEffect, useState, ReactNode } from "react";
import { User, Profile } from "@/types/f1";
import { LoginRequest, RegisterRequest } from "@/types/auth";
import { f1Api, ApiError } from "@/services/f1-api";

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
  const [user, setUser] = useState<User | null>(null);
  const [profile, setProfile] = useState<Profile | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    async function checkSession() {
      try {
        const response = await f1Api.getMe();
        
        if (response.user && response.user.id) {
            setUser(response.user);
            setProfile(response.profile);
        }
      } catch (err) {
        // Not logged in or error
        console.debug("No active session");
      } finally {
        setIsLoading(false);
      }
    }

    checkSession();
  }, []);

  const login = async (request: LoginRequest) => {
    setError(null);
    try {
      const response = await f1Api.login(request);
      setUser(response.user);
      setProfile(response.profile);
    } catch (err) {
      const message = err instanceof ApiError ? err.message : "Login failed";
      setError(message);
      throw err;
    }
  };

  const register = async (request: RegisterRequest) => {
    setError(null);
    try {
      const response = await f1Api.register(request);
      setUser(response.user);
      setProfile(response.profile);
    } catch (err) {
      const message = err instanceof ApiError ? err.message : "Registration failed";
      setError(message);
      throw err;
    }
  };

  const logout = async () => {
    try {
      await f1Api.logout();
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
