"use client";

import { useAuth } from "@/hooks/useAuth";
import { useRouter } from "next/navigation";
import { useEffect, useState } from "react";
import { LogOut, User as UserIcon, Mail, Calendar, Trophy, Copy, Check, UserPlus, Users, ChevronRight, ChevronLeft } from "lucide-react";
import { MainNav } from "@/components/ui/main-nav";
import { f1Api } from "@/services/f1-api";
import type { FriendRequest, LeaderboardEntry } from "@/types/f1";
import { Skeleton } from "@/components/ui/Skeleton";

export default function ProfilePage() {
  const { user, profile, isAuthenticated, isLoading, logout } = useAuth();
  const [selectedYear, setSelectedYear] = useState(2026);
  const [loadingProfile, setLoadingProfile] = useState(true);
  const [copied, setCopied] = useState(false);
  
  // Friends state
  const [pendingRequests, setPendingRequests] = useState<FriendRequest[]>([]);
  const [friendIdentifier, setFriendIdentifier] = useState("");
  const [sendingRequest, setSendingRequest] = useState(false);
  const [friendError, setFriendError] = useState("");
  const [friendSuccess, setFriendSuccess] = useState("");

  // Leaderboard state
  const [leaderboard, setLeaderboard] = useState<LeaderboardEntry[]>([]);
  const currentYear = new Date().getFullYear();
  const [loadingLeaderboard, setLoadingLeaderboard] = useState(false);
  const [errorFriends, setErrorFriends] = useState<string | null>(null);
  const [errorLeaderboard, setErrorLeaderboard] = useState<string | null>(null);

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
      setLoadingProfile(false);
      fetchFriendsData();
    }
  }, [isAuthenticated, user]);

  useEffect(() => {
    if (isAuthenticated) {
      fetchLeaderboard();
    }
  }, [selectedYear, isAuthenticated]);


  const fetchFriendsData = async () => {
    setErrorFriends(null);
    try {
      const requests = await f1Api.getPendingRequests();
      setPendingRequests(requests);
    } catch (err) {
      console.error("Failed to fetch friend requests:", err);
      setErrorFriends("Failed to load friend requests");
    }
  };

  const fetchLeaderboard = async () => {
    setLoadingLeaderboard(true);
    setErrorLeaderboard(null);
    try {
      const data = await f1Api.getLeaderboard(selectedYear);
      setLeaderboard(data);
    } catch (err) {
      console.error("Failed to fetch leaderboard:", err);
      setErrorLeaderboard("Failed to load leaderboard");
    } finally {
      setLoadingLeaderboard(false);
    }
  };

  const handleSendRequest = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!friendIdentifier) return;

    setSendingRequest(true);
    setFriendError("");
    setFriendSuccess("");

    try {
      await f1Api.sendFriendRequest(friendIdentifier);
      setFriendSuccess("Request sent successfully!");
      setFriendIdentifier("");
    } catch (err: unknown) {
      const message = err instanceof Error ? err.message : "Failed to send request";
      setFriendError(message);
    } finally {
      setSendingRequest(false);
    }
  };

  const handleRequestAction = async (requestId: string, action: "accept" | "deny") => {
    try {
      await f1Api.handleFriendRequest(requestId, action);
      await fetchFriendsData();
      if (action === "accept") {
        await fetchLeaderboard();
      }
    } catch (err) {
      console.error(`Failed to ${action} request:`, err);
    }
  };

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
          <main className="bg-slate-900/50 border border-slate-800 rounded-2xl overflow-hidden shadow-xl">
            <div className="p-8 space-y-6">
              <div className="flex items-center gap-6">
                <div className="h-24 w-24 rounded-full bg-red-600 flex items-center justify-center text-3xl font-bold border-4 border-slate-800 shadow-lg">
                  <UserIcon className="h-12 w-12" />
                </div>
                <div>
                  <h2 className="text-2xl font-bold">{profile.displayName}</h2>
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
                    <p>{new Date(user.createdAt).toLocaleDateString()}</p>
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

          {/* Friends Section */}
          <section className="space-y-6">
            <div className="flex items-center gap-2">
              <Users className="h-6 w-6 text-red-500" />
              <h3 className="text-xl font-bold">Friends</h3>
            </div>

            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
              {/* Add Friend Panel */}
              <div className="bg-slate-900/50 border border-slate-800 rounded-2xl p-6 space-y-4">
                <div className="flex items-center gap-2 text-sm font-semibold text-slate-400 uppercase tracking-wider">
                  <UserPlus className="h-4 w-4" />
                  Add a Friend
                </div>
                <form onSubmit={handleSendRequest} className="space-y-3">
                  <input
                    type="text"
                    placeholder="Email or Account ID"
                    value={friendIdentifier}
                    onChange={(e) => setFriendIdentifier(e.target.value)}
                    className="w-full bg-slate-950 border border-slate-800 rounded-lg px-4 py-2 text-sm focus:outline-none focus:border-red-500 transition-colors"
                  />
                  <button
                    type="submit"
                    disabled={sendingRequest || !friendIdentifier}
                    className="w-full bg-red-600 hover:bg-red-700 disabled:opacity-50 disabled:hover:bg-red-600 text-white font-semibold py-2 rounded-lg transition-colors text-sm"
                  >
                    {sendingRequest ? "Sending..." : "Send Invitation"}
                  </button>
                  {friendError && <p className="text-xs text-red-500">{friendError}</p>}
                  {friendSuccess && <p className="text-xs text-green-500">{friendSuccess}</p>}
                </form>
              </div>

              {/* Pending Requests Table */}
              <div className="bg-slate-900/50 border border-slate-800 rounded-2xl p-6 space-y-4">
                <div className="text-sm font-semibold text-slate-400 uppercase tracking-wider">
                  Incoming Requests
                </div>
                {errorFriends ? (
                  <div className="h-24 flex flex-col items-center justify-center rounded-xl bg-red-950/10 border border-dashed border-red-900/50">
                    <p className="text-red-500 text-xs">{errorFriends}</p>
                    <button 
                      onClick={fetchFriendsData}
                      className="text-[10px] text-slate-400 hover:underline mt-2"
                    >
                      Try Again
                    </button>
                  </div>
                ) : pendingRequests.length > 0 ? (
                  <div className="space-y-3">
                    {pendingRequests.map((req) => (
                      <div key={req.id} className="flex items-center justify-between p-3 rounded-xl bg-slate-950/50 border border-slate-800/50">
                        <div className="min-w-0">
                          <p className="font-bold text-sm truncate">{req.senderName}</p>
                          <p className="text-xs text-slate-500 truncate">{req.senderEmail}</p>
                        </div>
                        <div className="flex items-center gap-2 shrink-0">
                          <button
                            onClick={() => handleRequestAction(req.id, "accept")}
                            className="px-3 py-1 bg-green-600/20 hover:bg-green-600 text-green-500 hover:text-white rounded-lg text-xs font-bold transition-all"
                          >
                            Accept
                          </button>
                          <button
                            onClick={() => handleRequestAction(req.id, "deny")}
                            className="px-3 py-1 bg-red-600/20 hover:bg-red-600 text-red-500 hover:text-white rounded-lg text-xs font-bold transition-all"
                          >
                            Deny
                          </button>
                        </div>
                      </div>
                    ))}
                  </div>
                ) : (
                  <div className="h-24 flex items-center justify-center rounded-xl bg-slate-950/20 border border-dashed border-slate-800/50">
                    <p className="text-slate-600 text-xs text-center px-4">No pending invitations at the moment</p>
                  </div>
                )}
              </div>
            </div>
          </section>

          {/* Leaderboard Section */}
          <section className="space-y-6 pb-12">
            <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
              <div className="flex items-center gap-2">
                <Trophy className="h-6 w-6 text-yellow-500" />
                <h3 className="text-xl font-bold">Season Leaderboard</h3>
              </div>

              <div className="flex items-center gap-3 bg-slate-900/50 border border-slate-800 rounded-xl p-1 shrink-0">
                <button
                  onClick={() => setSelectedYear(selectedYear - 1)}
                  disabled={selectedYear <= 2025}
                  className="p-2 hover:bg-slate-800 rounded-lg transition-colors disabled:opacity-30 disabled:cursor-not-allowed"
                >
                  <ChevronLeft className="h-4 w-4" />
                </button>
                <span className="font-bold px-2">{selectedYear}</span>
                <button
                  onClick={() => setSelectedYear(selectedYear + 1)}
                  disabled={selectedYear >= currentYear}
                  className="p-2 hover:bg-slate-800 rounded-lg transition-colors disabled:opacity-30 disabled:cursor-not-allowed"
                >
                  <ChevronRight className="h-4 w-4" />
                </button>
              </div>
            </div>

            <div className="bg-slate-900/50 border border-slate-800 rounded-2xl overflow-hidden shadow-xl min-h-[300px]">
              <table className="w-full text-left">
                <thead>
                  <tr className="bg-slate-950/50 border-b border-slate-800">
                    <th className="px-6 py-4 text-xs font-bold text-slate-500 uppercase tracking-wider w-20">Pos</th>
                    <th className="px-6 py-4 text-xs font-bold text-slate-500 uppercase tracking-wider">User</th>
                    <th className="px-6 py-4 text-xs font-bold text-slate-500 uppercase tracking-wider text-right">Points</th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-slate-800/50">
                  {loadingLeaderboard ? (
                    [1, 2, 3, 4, 5].map(i => (
                      <tr key={i}>
                        <td colSpan={3} className="px-6 py-4">
                          <Skeleton className="h-10 w-full" />
                        </td>
                      </tr>
                    ))
                  ) : errorLeaderboard ? (
                    <tr>
                      <td colSpan={3} className="px-6 py-12 text-center">
                        <p className="text-red-500 text-sm mb-2">{errorLeaderboard}</p>
                        <button 
                          onClick={fetchLeaderboard}
                          className="text-xs text-slate-400 hover:underline"
                        >
                          Try Again
                        </button>
                      </td>
                    </tr>
                  ) : leaderboard.length > 0 ? (
                    leaderboard.map((entry) => (
                      <tr
                        key={entry.userId}
                        className={`transition-colors hover:bg-slate-800/30 ${entry.userId === user.id ? "bg-red-500/5" : ""}`}
                      >
                        <td className="px-6 py-4 font-mono text-sm">
                          <span className={`
                          ${entry.position === 1 ? "text-yellow-500 font-bold" : ""}
                          ${entry.position === 2 ? "text-slate-300 font-bold" : ""}
                          ${entry.position === 3 ? "text-amber-600 font-bold" : ""}
                          ${entry.position > 3 ? "text-slate-500" : ""}
                        `}>
                            #{entry.position}
                          </span>
                        </td>
                        <td className="px-6 py-4">
                          <div className="flex items-center gap-4">
                            <div className={`h-8 w-8 rounded-full flex items-center justify-center text-xs font-bold border-2 shrink-0 ${entry.userId === user.id ? "bg-red-600 border-red-400" : "bg-slate-800 border-slate-700"}`}>
                              <UserIcon className="h-4 w-4" />
                            </div>
                            <span className={`font-medium ${entry.userId === user.id ? "text-white" : "text-slate-300"}`}>
                              {entry.displayName} {entry.userId === user.id && <span className="text-[10px] bg-red-600 text-white px-1.5 py-0.5 rounded ml-2 font-bold inline-block align-middle">YOU</span>}
                            </span>
                          </div>
                        </td>
                        <td className="px-6 py-4 text-right">
                          <span className={`font-bold ${entry.userId === user.id ? "text-red-500" : "text-white"}`}>
                            {entry.score}
                          </span>
                        </td>
                      </tr>
                    ))
                  ) : (
                    <tr>
                      <td colSpan={3} className="px-6 py-12 text-center text-slate-500 text-sm italic">
                        No scoring data available for this season
                      </td>
                    </tr>
                  )}
                </tbody>
              </table>
            </div>
          </section>
        </div>
      </div>
    </div>
  );
}
