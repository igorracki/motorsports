import { User as UserIcon, Mail, Calendar, Copy, Check, LogOut } from "lucide-react";
import { useState } from "react";
import type { User, Profile } from "@/types/f1";

interface ProfileHeaderProps {
  user: User;
  profile: Profile;
  onLogout: () => void;
}

export function ProfileHeader({ user, profile, onLogout }: ProfileHeaderProps) {
  const [copied, setCopied] = useState(false);

  const handleCopyId = () => {
    navigator.clipboard.writeText(user.id)
      .then(() => {
        setCopied(true);
        setTimeout(() => setCopied(false), 2000);
      });
  };

  return (
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
                  <Copy className="h-3.5 w-3.5 text-slate-400" />
                )}
              </button>
            </div>
          </div>
        </div>
      </div>

      <div className="bg-slate-800/50 p-6 flex justify-end">
        <button
          onClick={onLogout}
          className="flex items-center gap-2 px-4 py-2 bg-slate-950 hover:bg-red-950 text-white hover:text-red-400 rounded-lg border border-slate-700 hover:border-red-900 transition-all duration-200"
        >
          <LogOut className="h-4 w-4" />
          Sign Out
        </button>
      </div>
    </main>
  );
}
