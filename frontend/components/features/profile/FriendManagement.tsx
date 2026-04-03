import { UserPlus, Users } from "lucide-react";
import { useState } from "react";
import type { FriendRequest } from "@/types/f1";
import { useAsync } from "@/hooks/useAsync";
import { useApi } from "@/components/providers/api-provider";

interface FriendManagementProps {
  pendingRequests: FriendRequest[];
  onAction: (requestId: string, action: "accept" | "deny") => Promise<void>;
  onRefresh: () => void;
  error?: string | null;
}

export function FriendManagement({ pendingRequests, onAction, onRefresh, error }: FriendManagementProps) {
  const { friendRepo } = useApi();
  const [friendIdentifier, setFriendIdentifier] = useState("");
  const [successMsg, setSuccessMsg] = useState("");

  const { execute: sendRequest, loading: sending, error: sendError } = useAsync(
    async (id: string) => {
      await friendRepo.sendFriendRequest(id);
    },
    {
      onSuccess: () => {
        setSuccessMsg("Request sent successfully!");
        setFriendIdentifier("");
        setTimeout(() => setSuccessMsg(""), 3000);
      }
    }
  );

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (friendIdentifier) {
      sendRequest(friendIdentifier);
    }
  };

  return (
    <section className="space-y-6">
      <div className="flex items-center gap-2">
        <Users className="h-6 w-6 text-red-500" />
        <h3 className="text-xl font-bold">Friends</h3>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
        <div className="bg-slate-900/50 border border-slate-800 rounded-2xl p-6 space-y-4">
          <div className="flex items-center gap-2 text-sm font-semibold text-slate-400 uppercase tracking-wider">
            <UserPlus className="h-4 w-4" />
            Add a Friend
          </div>
          <form onSubmit={handleSubmit} className="space-y-3">
            <input
              type="text"
              placeholder="Email or Account ID"
              value={friendIdentifier}
              onChange={(e) => setFriendIdentifier(e.target.value)}
              className="w-full bg-slate-950 border border-slate-800 rounded-lg px-4 py-2 text-sm focus:outline-none focus:border-red-500 transition-colors"
            />
            <button
              type="submit"
              disabled={sending || !friendIdentifier}
              className="w-full bg-red-600 hover:bg-red-700 disabled:opacity-50 disabled:hover:bg-red-600 text-white font-semibold py-2 rounded-lg transition-colors text-sm"
            >
              {sending ? "Sending..." : "Send Invitation"}
            </button>
            {sendError && <p className="text-xs text-red-500">{sendError}</p>}
            {successMsg && <p className="text-xs text-green-500">{successMsg}</p>}
          </form>
        </div>

        <div className="bg-slate-900/50 border border-slate-800 rounded-2xl p-6 space-y-4">
          <div className="text-sm font-semibold text-slate-400 uppercase tracking-wider">
            Incoming Requests
          </div>
          {error ? (
            <div className="h-24 flex flex-col items-center justify-center rounded-xl bg-red-950/10 border border-dashed border-red-900/50">
              <p className="text-red-500 text-xs">{error}</p>
              <button 
                onClick={onRefresh}
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
                      onClick={() => onAction(req.id, "accept")}
                      className="px-3 py-1 bg-green-600/20 hover:bg-green-600 text-green-500 hover:text-white rounded-lg text-xs font-bold transition-all"
                    >
                      Accept
                    </button>
                    <button
                      onClick={() => onAction(req.id, "deny")}
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
  );
}
