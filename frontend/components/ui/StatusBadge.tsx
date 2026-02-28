import { cn } from "@/lib/utils";

interface StatusBadgeProps {
  status: "completed" | "ongoing" | "upcoming";
  className?: string;
}

export function StatusBadge({ status, className }: StatusBadgeProps) {
  if (status === "completed") return null;

  if (status === "ongoing") {
    return (
      <div className={cn("flex items-center gap-1.5 rounded-md bg-primary/20 px-2 py-0.5 text-xs font-semibold text-primary", className)}>
        <span className="relative flex h-2 w-2">
          <span className="absolute inline-flex h-full w-full animate-ping rounded-full bg-primary opacity-75" />
          <span className="relative inline-flex h-2 w-2 rounded-full bg-primary" />
        </span>
        ONGOING
      </div>
    );
  }

  return (
    <div className={cn("rounded-md bg-muted px-2 py-0.5 text-xs font-medium text-muted-foreground", className)}>
      UPCOMING
    </div>
  );
}
