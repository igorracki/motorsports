"use client";

import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { RegisterRequest, RegisterRequestSchema } from "@/types/auth";
import { useAuth } from "@/hooks/useAuth";
import { useAsync } from "@/hooks/useAsync";
import { AlertCircle, Loader2, Mail, Lock, User } from "lucide-react";
import Link from "next/link";

import { Input } from "@/components/ui/Input";
import { Button } from "@/components/ui/Button";

interface RegisterFormProps {
  onSuccess?: () => void;
}

export function RegisterForm({ onSuccess }: RegisterFormProps) {
  const { register: registerUser } = useAuth();

  const { execute: performRegister, loading: isSubmitting, error: serverError } = useAsync(
    async (data: RegisterRequest) => {
      await registerUser(data);
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
  } = useForm<RegisterRequest>({
    resolver: zodResolver(RegisterRequestSchema),
    defaultValues: {
      email: "",
      password: "",
      display_name: "",
    },
  });

  const onSubmit = async (data: RegisterRequest) => {
    try {
      await performRegister(data);
    } catch {
      // Error handled by useAsync and displayed via serverError
    }
  };

  return (
    <div className="w-full max-w-md space-y-6 p-6 bg-slate-900/50 backdrop-blur-sm rounded-xl border border-slate-800 shadow-xl">
      <div className="space-y-2 text-center">
        <h1 className="text-3xl font-bold tracking-tighter text-white">Create Account</h1>
        <p className="text-slate-400">Join the platform to start predicting</p>
      </div>

      <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
        {serverError && (
          <div className="p-3 rounded-lg bg-red-500/10 border border-red-500/20 flex items-center gap-3 text-red-500 text-sm">
            <AlertCircle className="h-4 w-4 shrink-0" />
            <p>{serverError}</p>
          </div>
        )}

        <div className="space-y-2">
          <label className="text-sm font-medium text-slate-300" htmlFor="display_name">
            Display Name
          </label>
          <Input
            {...register("display_name")}
            id="display_name"
            type="text"
            placeholder="jimmy_v"
            icon={<User className="h-4 w-4" />}
          />
          {errors.display_name && <p className="text-xs text-red-500">{errors.display_name.message}</p>}
        </div>

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
          <label className="text-sm font-medium text-slate-300" htmlFor="password">
            Password
          </label>
          <Input
            {...register("password")}
            id="password"
            type="password"
            placeholder="••••••••"
            icon={<Lock className="h-4 w-4" />}
          />
          {errors.password && <p className="text-xs text-red-500">{errors.password.message}</p>}
        </div>

        <Button
          type="submit"
          disabled={isSubmitting}
          fullWidth
          className="mt-2"
        >
          {isSubmitting ? <Loader2 className="h-4 w-4 animate-spin" /> : null}
          {isSubmitting ? "Creating account..." : "Register"}
        </Button>
      </form>

      <div className="text-center text-sm text-slate-400">
        Already have an account?{" "}
        <Link href="/login" className="text-red-500 hover:text-red-400 font-medium">
          Login
        </Link>
      </div>
    </div>
  );
}
