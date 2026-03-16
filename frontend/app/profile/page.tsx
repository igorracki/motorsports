"use client";

import { useAuth } from "@/hooks/useAuth";
import { useRouter } from "next/navigation";
import { useEffect, useState } from "react";
import { LogOut, User as UserIcon, Mail, Calendar, Trophy, Copy, Check } from "lucide-react";
import { MainNav } from "@/components/main-nav";
import { f1Api } from "@/services/f1-api";
import type { UserProfileResponse, UserScore } from "@/types/f1";
import { Skeleton } from "@/components/ui/Skeleton";

export default function ProfilePage() {
  const { user, profile, isAuthenticated, isLoading, logout } = useAuth();
  const [fullProfile, setFullProfile] = useState<UserProfileResponse | null>(null);
  const [seasonScores, setSeasonScores] = useState<UserScore[]>([]);
  const [loadingProfile, setLoadingProfile] = useState(true);
  const [loadingScores, setLoadingScores] = useState(true);
  const [copied, setCopied] = useState(false);
  const router = useRouter();

  const handleCopyId = () => {
    if (!user?.id) return;
    navigator.clipboard.writeText(user.id);
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  };

  useEffect(() => {
    if (!isLoading && !isAuthenticated) {
      router.push("/login");
    }
  }, [isAuthenticated, isLoading, router]);

  useEffect(() => {
    if (isAuthenticated && user) {
      setLoadingProfile(true);
      f1Api.getUserProfile(user.id)
        .then(setFullProfile)
        .catch(err => console.error("Failed to fetch full profile:", err))
        .finally(() => setLoadingProfile(false));

      setLoadingScores(true);
      f1Api.getSeasonScores(user.id)
        .then(setSeasonScores)
        .catch(err => console.error("Failed to fetch season scores:", err))
        .finally(() => setLoadingScores(false));
    }
  }, [isAuthenticated, user]);

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
    <div className="min-h-screen bg-background text-white pb-20">
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
                <div className="flex items-center justify-between gap-2">
                  <div className="flex items-center gap-2 overflow-hidden">
                    <UserIcon className="h-4 w-4 text-red-500 shrink-0" />
                    <p className="font-mono text-xs truncate" title={user.id}>{user.id}</p>
                  </div>
                  <button 
                    onClick={handleCopyId}
                    className="p-1.5 hover:bg-slate-800 rounded-lg transition-colors shrink-0"
                    title="Copy ID"
                  >
                    {copied ? (
                      <Check className="h-3.5 w-3.5 text-green-500" />
                    ) : (
                      <Copy className="h-3.5 w-3.5 text-slate-400 group-hover:text-white" />
                    )}
                  </button>
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

        <section className="space-y-6">
          <div className="flex items-center justify-between">
            <h3 className="text-xl font-bold">Prediction Dashboard</h3>
            {loadingScores && <Skeleton className="h-4 w-24" />}
          </div>
          
          {loadingScores ? (
            <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
              {[1, 2].map(i => (
                <Skeleton key={i} className="h-32 rounded-2xl" />
              ))}
            </div>
          ) : seasonScores.length > 0 ? (
            <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
              {seasonScores.sort((a, b) => (b.season || 0) - (a.season || 0)).map((score) => (
                <div 
                  key={score.season}
                  className="group relative overflow-hidden rounded-2xl border border-slate-800 bg-slate-950/50 p-6 transition-all duration-300 hover:border-red-500/50 hover:shadow-lg hover:shadow-red-500/5"
                >
                  <div className="relative z-10 flex items-center justify-between">
                    <div>
                      <p className="text-sm font-semibold uppercase tracking-wider text-slate-500 group-hover:text-red-400/80 transition-colors">
                        {score.season} Season
                      </p>
                      <div className="mt-1 flex items-baseline gap-2">
                        <span className="text-3xl font-bold text-white transition-transform group-hover:scale-105 inline-block">
                          {score.value}
                        </span>
                        <span className="text-sm font-medium text-slate-400">Points</span>
                      </div>
                    </div>
                    <div className="rounded-xl bg-red-500/10 p-3 text-red-500 transition-colors group-hover:bg-red-500/20">
                      <Trophy className="h-6 w-6" />
                    </div>
                  </div>
                  {/* Decorative background element */}
                  <div className="absolute -right-4 -top-4 h-24 w-24 rounded-full bg-red-500/5 blur-2xl group-hover:bg-red-500/10 transition-colors" />
                </div>
              ))}
            </div>
          ) : (
            <div className="h-32 rounded-xl bg-slate-950/50 border border-dashed border-slate-800 flex items-center justify-center animate-in fade-in duration-500">
              <p className="text-slate-600 text-sm">No recent prediction history found</p>
            </div>
          )}
        </section>
      </div>
    </div>
  </div>
  );
}
