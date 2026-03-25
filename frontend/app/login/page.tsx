"use client";

import { LoginForm } from "@/components/features/auth/LoginForm";
import { useRouter } from "next/navigation";
import { useAuth } from "@/hooks/useAuth";
import { useEffect } from "react";
import { MainNav } from "@/components/ui/main-nav";
import { Footer } from "@/components/ui/Footer";

export default function LoginPage() {
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
        <LoginForm onSuccess={() => router.push("/")} />
      </div>
      <Footer />
    </div>
  );
}
