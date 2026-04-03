import { useState, useCallback, useEffect } from "react";
import { useApi } from "@/components/providers/api-provider";
import { useAsync } from "./useAsync";
import type { FriendRequest, LeaderboardEntry } from "@/types/f1";

export function useProfileController(selectedYear: number) {
  const { friendRepo } = useApi();
  const [pendingRequests, setPendingRequests] = useState<FriendRequest[]>([]);
  const [leaderboard, setLeaderboard] = useState<LeaderboardEntry[]>([]);

  const { execute: fetchRequests, loading: loadingRequests, error: errorFriends } = useAsync(
    async () => {
      const requests = await friendRepo.getPendingRequests();
      setPendingRequests(requests);
      return requests;
    }
  );

  const { execute: fetchLeaderboard, loading: loadingLeaderboard, error: errorLeaderboard } = useAsync(
    async (year: number) => {
      const data = await friendRepo.getLeaderboard(year);
      setLeaderboard(data);
      return data;
    }
  );

  const handleRequestAction = useCallback(async (requestId: string, action: "accept" | "deny") => {
    try {
      await friendRepo.handleFriendRequest(requestId, action);
      await fetchRequests();
      if (action === "accept") {
        await fetchLeaderboard(selectedYear);
      }
    } catch {
      // Error handled by the trigger
    }
  }, [fetchRequests, fetchLeaderboard, selectedYear, friendRepo]);

  useEffect(() => {
    fetchRequests();
  }, [fetchRequests]);

  useEffect(() => {
    fetchLeaderboard(selectedYear);
  }, [selectedYear, fetchLeaderboard]);

  return {
    pendingRequests,
    leaderboard,
    loading: {
      requests: loadingRequests,
      leaderboard: loadingLeaderboard
    },
    errors: {
      friends: errorFriends,
      leaderboard: errorLeaderboard
    },
    actions: {
      fetchRequests,
      fetchLeaderboard,
      handleRequestAction
    }
  };
}
