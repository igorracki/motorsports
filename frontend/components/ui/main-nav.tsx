"use client";

import Link from "next/link";
import { Home, User, Settings, LogIn } from "lucide-react";
import { useAuth } from "@/hooks/useAuth";
import { useState, useRef, useEffect } from "react";
import { useSettings } from "@/hooks/useSettings";
import { Switch } from "@/components/ui/Switch";
import { cn } from "@/lib/utils";

export function MainNav() {
  const { isAuthenticated, profile } = useAuth();
  const { useBrowserTime, setUseBrowserTime } = useSettings();
  const [isOpen, setIsOpen] = useState(false);
  const containerRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    function handleClickOutside(event: MouseEvent) {
      if (containerRef.current && !containerRef.current.contains(event.target as Node)) {
        setIsOpen(false);
      }
    }
    document.addEventListener("mousedown", handleClickOutside);
    return () => document.removeEventListener("mousedown", handleClickOutside);
  }, []);

  return (
    <header className="sticky top-0 z-[100] w-full border-b border-border/40 bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
      <div className="container mx-auto flex h-14 items-center justify-between px-4 md:px-6">
        <Link href="/" className="flex items-center gap-2">
          <span className="hidden text-lg font-semibold tracking-tight text-foreground sm:inline-block">
            Motorsports - F1
          </span>
        </Link>

        <nav className="flex items-center gap-1 md:gap-2">
          <Link
            href="/"
            className="flex items-center gap-2 rounded-lg px-3 py-2 text-sm font-medium text-foreground/80 transition-colors hover:bg-secondary hover:text-foreground"
          >
            <Home className="h-4 w-4" />
            <span className="hidden sm:inline">Home</span>
          </Link>
          
          {isAuthenticated ? (
            <Link
              href="/profile"
              className="flex items-center gap-2 rounded-lg px-3 py-2 text-sm font-medium text-foreground/80 transition-colors hover:bg-secondary hover:text-foreground"
            >
              <User className="h-4 w-4" />
              <span className="hidden sm:inline">{profile?.displayName || "Profile"}</span>
            </Link>
          ) : (
            <Link
              href="/login"
              className="flex items-center gap-2 rounded-lg px-3 py-2 text-sm font-medium text-foreground/80 transition-colors hover:bg-secondary hover:text-foreground"
            >
              <LogIn className="h-4 w-4" />
              <span className="hidden sm:inline">Sign In</span>
            </Link>
          )}
          
          <div 
            className="relative"
            ref={containerRef}
          >
            <div
              className={cn(
                "flex items-center gap-2 rounded-lg px-3 py-2 text-sm font-medium text-foreground/80 transition-colors hover:bg-secondary hover:text-foreground cursor-pointer",
                isOpen && "bg-secondary text-foreground"
              )}
              onClick={() => setIsOpen(!isOpen)}
            >
              <Settings className="h-4 w-4" />
              <span className="hidden sm:inline">Settings</span>
            </div>

            {isOpen && (
              <div className="absolute right-0 top-full mt-1 w-64 origin-top-right rounded-xl border border-border/50 bg-card p-4 shadow-xl animate-in fade-in zoom-in-95 duration-200 z-50">
                <div className="flex items-center justify-between gap-4">
                  <div className="space-y-0.5">
                    <p className="text-sm font-semibold">Browser Local Times</p>
                    <p className="text-[10px] text-muted-foreground leading-tight">
                      Convert session times to your local timezone
                    </p>
                  </div>
                  <Switch 
                    checked={useBrowserTime}
                    onCheckedChange={setUseBrowserTime}
                  />
                </div>
              </div>
            )}
          </div>
        </nav>
      </div>
    </header>
  );
}
