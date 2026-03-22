"use client";

import React, { useState } from "react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { LoginRequest, LoginRequestSchema } from "@/types/auth";
import { useAuth } from "@/hooks/useAuth";
import { AlertCircle, Loader2, Mail, Lock } from "lucide-react";
import { cn } from "@/lib/utils";
import Link from "next/link";

interface LoginFormProps {
  onSuccess?: () => void;
  className?: string;
  hideContainer?: boolean;
}

export function LoginForm({ onSuccess, className, hideContainer }: LoginFormProps) {
  const { login } = useAuth();
  const [serverError, setServerError] = useState<string | null>(null);
  const [isSubmitting, setIsSubmitting] = useState(false);

  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<LoginRequest>({
    resolver: zodResolver(LoginRequestSchema),
    defaultValues: {
      email: "",
      password: "",
      remember_me: true,
    },
  });

  const onSubmit = async (data: LoginRequest) => {
    setIsSubmitting(true);
    setServerError(null);
    try {
      await login(data);
      if (onSuccess) onSuccess();
    } catch (err: unknown) {
      const message = err instanceof Error ? err.message : "Invalid email or password";
      setServerError(message);
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <div className={cn(
      "w-full space-y-6",
      !hideContainer && "max-w-md p-6 bg-slate-900/50 backdrop-blur-sm rounded-xl border border-slate-800 shadow-xl",
      className
    )}>
      <div className="space-y-2 text-center">
        <h1 className="text-3xl font-bold tracking-tighter text-white">Login</h1>
        <p className="text-slate-400">Enter your credentials to access your account</p>
      </div>

      <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
        {serverError && (
          <div className="p-3 rounded-lg bg-red-500/10 border border-red-500/20 flex items-center gap-3 text-red-500 text-sm">
            <AlertCircle className="h-4 w-4 shrink-0" />
            <p>{serverError}</p>
          </div>
        )}

        <div className="space-y-2">
          <label className="text-sm font-medium text-slate-300" htmlFor="email">
            Email
          </label>
          <div className="relative">
            <Mail className="absolute left-3 top-3 h-4 w-4 text-slate-500" />
            <input
              {...register("email")}
              id="email"
              type="email"
              placeholder="name@example.com"
              className="w-full bg-slate-950 border border-slate-800 rounded-lg py-2 pl-10 pr-4 text-white placeholder:text-slate-600 focus:outline-none focus:ring-2 focus:ring-red-500/50 transition-all"
            />
          </div>
          {errors.email && <p className="text-xs text-red-500">{errors.email.message}</p>}
        </div>

        <div className="space-y-2">
          <div className="flex items-center justify-between">
            <label className="text-sm font-medium text-slate-300" htmlFor="password">
              Password
            </label>
          </div>
          <div className="relative">
            <Lock className="absolute left-3 top-3 h-4 w-4 text-slate-500" />
            <input
              {...register("password")}
              id="password"
              type="password"
              className="w-full bg-slate-950 border border-slate-800 rounded-lg py-2 pl-10 pr-4 text-white placeholder:text-slate-600 focus:outline-none focus:ring-2 focus:ring-red-500/50 transition-all"
            />
          </div>
          {errors.password && <p className="text-xs text-red-500">{errors.password.message}</p>}
        </div>

        <div className="flex items-center space-x-2">
          <input
            {...register("remember_me")}
            type="checkbox"
            id="remember_me"
            className="h-4 w-4 rounded border-slate-800 bg-slate-950 text-red-600 focus:ring-red-500/50"
          />
          <label htmlFor="remember_me" className="text-sm text-slate-400">
            Remember me
          </label>
        </div>

        <button
          type="submit"
          disabled={isSubmitting}
          className="w-full bg-red-600 hover:bg-red-700 text-white font-semibold py-2 rounded-lg transition-colors flex items-center justify-center gap-2 disabled:opacity-50 disabled:cursor-not-allowed"
        >
          {isSubmitting ? <Loader2 className="h-4 w-4 animate-spin" /> : null}
          {isSubmitting ? "Logging in..." : "Login"}
        </button>
      </form>

      <div className="text-center text-sm text-slate-400">
        Don't have an account?{" "}
        <Link href="/register" className="text-red-500 hover:text-red-400 font-medium">
          Register
        </Link>
      </div>
    </div>
  );
}
