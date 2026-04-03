"use client";

import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { LoginRequest, LoginRequestSchema } from "@/types/auth";
import { useAuth } from "@/hooks/useAuth";
import { useAsync } from "@/hooks/useAsync";
import { AlertCircle, Loader2, Mail, Lock } from "lucide-react";
import { cn } from "@/lib/utils";
import Link from "next/link";
import { Input } from "@/components/ui/Input";
import { Button } from "@/components/ui/Button";

export interface LoginFormProps {
  onSuccess?: () => void;
  className?: string;
  hideContainer?: boolean;
}

export function LoginForm({ onSuccess, className, hideContainer }: LoginFormProps) {
  const { login } = useAuth();

  const { execute: performLogin, loading: isSubmitting, error: serverError } = useAsync(
    async (data: LoginRequest) => {
      await login(data);
    },
    {
      onSuccess: () => {
        if (onSuccess) onSuccess();
      },
    }
  );

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
    try {
      await performLogin(data);
    } catch {
      // Error handled by useAsync and displayed via serverError
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
          <Input
            {...register("email")}
            id="email"
            type="email"
            placeholder="name@example.com"
            icon={<Mail className="h-4 w-4" />}
          />
          {errors.email && <p className="text-xs text-red-500">{errors.email.message}</p>}
        </div>

        <div className="space-y-2">
          <div className="flex items-center justify-between">
            <label className="text-sm font-medium text-slate-300" htmlFor="password">
              Password
            </label>
          </div>
          <Input
            {...register("password")}
            id="password"
            type="password"
            icon={<Lock className="h-4 w-4" />}
          />
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

        <Button
          type="submit"
          disabled={isSubmitting}
          fullWidth
        >
          {isSubmitting ? <Loader2 className="h-4 w-4 animate-spin" /> : null}
          {isSubmitting ? "Logging in..." : "Login"}
        </Button>
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

