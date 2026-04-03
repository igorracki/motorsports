"use client";

import { useAuth } from "@/hooks/useAuth";
import { useRouter } from "next/navigation";
import { useEffect, useState } from "react";
import { MainNav } from "@/components/ui/main-nav";
import { Footer } from "@/components/ui/Footer";
import { useProfileController } from "@/hooks/useProfileController";
import { ErrorBoundary } from "@/components/ui/error-boundary";

import { ProfileHeader } from "@/components/features/profile/ProfileHeader";
import { FriendManagement } from "@/components/features/profile/FriendManagement";
import { SeasonLeaderboard } from "@/components/features/profile/SeasonLeaderboard";

export default function ProfilePage() {
  const { user, profile, isAuthenticated, isLoading, logout } = useAuth();
  const [selectedYear, setSelectedYear] = useState(() => new Date().getFullYear());
  const currentYear = new Date().getFullYear();
  const router = useRouter();

  const {
    pendingRequests,
    leaderboard,
    loading,
    errors,
    actions
  } = useProfileController(selectedYear);

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

  const handleLogout = async () => {
    await logout();
    router.push("/login");
  };

  return (
    <div className="min-h-screen bg-background text-white flex flex-col">
      <MainNav />

      <div className="flex-1 p-4 md:p-8">
        <div className="max-w-2xl mx-auto space-y-8">
          <ErrorBoundary name="Profile Header">
            <ProfileHeader
              user={user}
              profile={profile}
              onLogout={handleLogout}
            />
          </ErrorBoundary>

          <ErrorBoundary name="Friend Management">
            <FriendManagement
              pendingRequests={pendingRequests}
              onAction={actions.handleRequestAction}
              onRefresh={actions.fetchRequests}
              error={errors.friends}
            />
          </ErrorBoundary>

          <ErrorBoundary name="Season Leaderboard">
            <SeasonLeaderboard
              leaderboard={leaderboard}
              selectedYear={selectedYear}
              currentYear={currentYear}
              loading={loading.leaderboard}
              error={errors.leaderboard}
              currentUserId={user.id}
              onYearChange={setSelectedYear}
              onRefresh={() => actions.fetchLeaderboard(selectedYear)}
            />
          </ErrorBoundary>
        </div>
      </div>
      <Footer />
    </div>
  );
}
