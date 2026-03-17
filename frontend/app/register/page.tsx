"use client";

import { RegisterForm } from "@/components/features/auth/RegisterForm";
import { useRouter } from "next/navigation";
import { useAuth } from "@/hooks/useAuth";
import { useEffect } from "react";
import { MainNav } from "@/components/ui/main-nav";

export default function RegisterPage() {
  const router = useRouter();
  const { isAuthenticated, isLoading } = useAuth();

  useEffect(() => {
    if (isAuthenticated) {
      router.push("/");
    }
  }, [isAuthenticated, router]);

  if (isLoading || isAuthenticated) {
    return null;
  }

  return (
    <div className="min-h-screen bg-background flex flex-col">
      <MainNav />
      <div className="flex-1 flex flex-col items-center justify-center p-4">
        <RegisterForm onSuccess={() => router.push("/")} />
      </div>
    </div>
  );
}
