import { Trophy, ChevronLeft, ChevronRight, User as UserIcon } from "lucide-react";
import { Skeleton } from "@/components/ui/Skeleton";
import type { LeaderboardEntry } from "@/types/f1";

interface SeasonLeaderboardProps {
  leaderboard: LeaderboardEntry[];
  selectedYear: number;
  currentYear: number;
  loading: boolean;
  error: string | null;
  currentUserId: string;
  onYearChange: (year: number) => void;
  onRefresh: () => void;
}

export function SeasonLeaderboard({
  leaderboard,
  selectedYear,
  currentYear,
  loading,
  error,
  currentUserId,
  onYearChange,
  onRefresh
}: SeasonLeaderboardProps) {
  return (
    <section className="space-y-6 pb-12">
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
        <div className="flex items-center gap-2">
          <Trophy className="h-6 w-6 text-yellow-500" />
          <h3 className="text-xl font-bold">Season Leaderboard</h3>
        </div>

        <div className="flex items-center gap-3 bg-slate-900/50 border border-slate-800 rounded-xl p-1 shrink-0">
          <button
            onClick={() => onYearChange(selectedYear - 1)}
            disabled={selectedYear <= 2025}
            className="p-2 hover:bg-slate-800 rounded-lg transition-colors disabled:opacity-30 disabled:cursor-not-allowed"
          >
            <ChevronLeft className="h-4 w-4" />
          </button>
          <span className="font-bold px-2">{selectedYear}</span>
          <button
            onClick={() => onYearChange(selectedYear + 1)}
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
            {loading ? (
              [1, 2, 3, 4, 5].map(i => (
                <tr key={i}>
                  <td colSpan={3} className="px-6 py-4">
                    <Skeleton className="h-10 w-full" />
                  </td>
                </tr>
              ))
            ) : error ? (
              <tr>
                <td colSpan={3} className="px-6 py-12 text-center">
                  <p className="text-red-500 text-sm mb-2">{error}</p>
                  <button 
                    onClick={onRefresh}
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
                  className={`transition-colors hover:bg-slate-800/30 ${entry.userId === currentUserId ? "bg-red-500/5" : ""}`}
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
                      <div className={`h-8 w-8 rounded-full flex items-center justify-center text-xs font-bold border-2 shrink-0 ${entry.userId === currentUserId ? "bg-red-600 border-red-400" : "bg-slate-800 border-slate-700"}`}>
                        <UserIcon className="h-4 w-4" />
                      </div>
                      <span className={`font-medium ${entry.userId === currentUserId ? "text-white" : "text-slate-300"}`}>
                        {entry.displayName} {entry.userId === currentUserId && <span className="text-[10px] bg-red-600 text-white px-1.5 py-0.5 rounded ml-2 font-bold inline-block align-middle">YOU</span>}
                      </span>
                    </div>
                  </td>
                  <td className="px-6 py-4 text-right">
                    <span className={`font-bold ${entry.userId === currentUserId ? "text-red-500" : "text-white"}`}>
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
  );
}
