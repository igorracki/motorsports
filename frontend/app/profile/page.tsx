"use client";

import { useAuth } from "@/hooks/useAuth";
import { useRouter } from "next/navigation";
import { useEffect } from "react";
import { LogOut, User as UserIcon, Mail, Calendar } from "lucide-react";
import { MainNav } from "@/components/main-nav";

export default function ProfilePage() {
  const { user, profile, isAuthenticated, isLoading, logout } = useAuth();
  const router = useRouter();

  useEffect(() => {
    if (!isLoading && !isAuthenticated) {
      router.push("/login");
    }
  }, [isAuthenticated, isLoading, router]);

  if (isLoading) {
    return (
      <div className="min-h-screen bg-background flex items-center justify-center">
        <div className="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-red-500"></div>
      </div>
    );
  }

  if (!isAuthenticated || !user || !profile) {
    return null;
  }

  return (
    <div className="min-h-screen bg-background text-white">
      <MainNav />
      
      <div className="p-4 md:p-8">
        <div className="max-w-2xl mx-auto space-y-8">
        <header className="space-y-2">
          <h1 className="text-4xl font-bold tracking-tight">Profile</h1>
          <p className="text-slate-400">Manage your account and view your stats</p>
        </header>

        <main className="bg-slate-900/50 border border-slate-800 rounded-2xl overflow-hidden shadow-xl">
          <div className="p-8 space-y-6">
            <div className="flex items-center gap-6">
              <div className="h-24 w-24 rounded-full bg-red-600 flex items-center justify-center text-3xl font-bold border-4 border-slate-800 shadow-lg">
                {profile.display_name.charAt(0).toUpperCase()}
              </div>
              <div>
                <h2 className="text-2xl font-bold">{profile.display_name}</h2>
                <p className="text-slate-400 flex items-center gap-2">
                  <Mail className="h-4 w-4" /> {user.email}
                </p>
              </div>
            </div>

            <div className="grid grid-cols-1 md:grid-cols-2 gap-4 pt-6 border-t border-slate-800">
              <div className="p-4 rounded-xl bg-slate-950/50 border border-slate-800/50 space-y-1">
                <p className="text-xs font-medium text-slate-500 uppercase tracking-wider">Member Since</p>
                <div className="flex items-center gap-2">
                  <Calendar className="h-4 w-4 text-red-500" />
                  <p>{new Date(user.created_at).toLocaleDateString()}</p>
                </div>
              </div>
              <div className="p-4 rounded-xl bg-slate-950/50 border border-slate-800/50 space-y-1">
                <p className="text-xs font-medium text-slate-500 uppercase tracking-wider">Account ID</p>
                <div className="flex items-center gap-2">
                  <UserIcon className="h-4 w-4 text-red-500" />
                  <p className="font-mono text-xs">{user.id.substring(0, 8)}...</p>
                </div>
              </div>
            </div>
          </div>

          <div className="bg-slate-800/50 p-6 flex justify-end">
            <button
              onClick={async () => {
                await logout();
                router.push("/login");
              }}
              className="flex items-center gap-2 px-4 py-2 bg-slate-950 hover:bg-red-950 text-white hover:text-red-400 rounded-lg border border-slate-700 hover:border-red-900 transition-all duration-200"
            >
              <LogOut className="h-4 w-4" />
              Sign Out
            </button>
          </div>
        </main>

        <section className="bg-slate-900/50 border border-slate-800 rounded-2xl p-8 space-y-4">
          <h3 className="text-xl font-bold">Prediction Dashboard</h3>
          <p className="text-slate-400 italic">Dashboard coming soon... Stay tuned for your performance metrics!</p>
          <div className="h-32 rounded-xl bg-slate-950/50 border border-dashed border-slate-800 flex items-center justify-center">
             <p className="text-slate-600 text-sm">No recent prediction history found</p>
          </div>
        </section>
      </div>
    </div>
  </div>
  );
}
