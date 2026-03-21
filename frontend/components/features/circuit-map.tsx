"use client";

import { useMemo, useId } from "react";
import type { CircuitLayoutPoint } from "@/types/f1";

interface CircuitMapProps {
  layout?: CircuitLayoutPoint[];
  rotation?: number;
  className?: string;
}

export function CircuitMap({ layout, rotation = 0, className }: CircuitMapProps) {
  const gradientId = useId();
  const glowId = useId();
  const result = useMemo(() => {
    if (!layout || layout.length === 0) return null;

    // 1. Apply rotation to raw points first to maximize space usage in the box
    const rad = (rotation * Math.PI) / 180;
    const rotated = layout.map((p) => ({
      x: p.x * Math.cos(rad) - p.y * Math.sin(rad),
      y: p.x * Math.sin(rad) + p.y * Math.cos(rad),
    }));

    // 2. Find bounds of rotated track
    let minX = Infinity,
      maxX = -Infinity,
      minY = Infinity,
      maxY = -Infinity;

    rotated.forEach((p) => {
      if (p.x < minX) minX = p.x;
      if (p.x > maxX) maxX = p.x;
      if (p.y < minY) minY = p.y;
      if (p.y > maxY) maxY = p.y;
    });

    const width = maxX - minX;
    const height = maxY - minY;

    // Guard against degenerate layouts (e.g., single point or collinear points)
    if (width === 0 && height === 0) return null;

    // 3. Scale and center to viewBox (300x260) with minimal padding
    const padding = 20;
    const targetW = 300 - padding * 2;
    const targetH = 260 - padding * 2;

    const scale = Math.min(
      width > 0 ? targetW / width : Infinity,
      height > 0 ? targetH / height : Infinity
    );

    // Center the track in the 300x260 area
    const offsetX = (300 - width * scale) / 2;
    const offsetY = (260 - height * scale) / 2;

    const points = rotated.map((p) => ({
      x: (p.x - minX) * scale + offsetX,
      y: (p.y - minY) * scale + offsetY,
    }));

    const d =
      points.map((p, i) => `${i === 0 ? "M" : "L"}${p.x},${p.y}`).join(" ") +
      " Z";

    return {
      d,
      firstPoint: points[0],
    };
  }, [layout, rotation]);

  return (
    <div className={className}>
      {result ? (
        <svg
          viewBox="0 0 300 260"
          className="h-full w-full"
          fill="none"
          xmlns="http://www.w3.org/2000/svg"
        >
          <defs>
            <linearGradient
              id={gradientId}
              x1="0%"
              y1="0%"
              x2="100%"
              y2="100%"
            >
              <stop
                offset="0%"
                stopColor="oklch(0.65 0.2 25)"
                stopOpacity="0.8"
              />
              <stop
                offset="100%"
                stopColor="oklch(0.75 0.15 85)"
                stopOpacity="0.6"
              />
            </linearGradient>
            <filter id={glowId}>
              <feGaussianBlur stdDeviation="3" result="coloredBlur" />
              <feMerge>
                <feMergeNode in="coloredBlur" />
                <feMergeNode in="SourceGraphic" />
              </feMerge>
            </filter>
          </defs>

          {/* Background glow */}
          <path
            d={result.d}
            stroke={`url(#${gradientId})`}
            strokeWidth="12"
            strokeLinecap="round"
            strokeLinejoin="round"
            opacity="0.2"
            filter={`url(#${glowId})`}
          />

          {/* Main track */}
          <path
            d={result.d}
            stroke={`url(#${gradientId})`}
            strokeWidth="6"
            strokeLinecap="round"
            strokeLinejoin="round"
            className="transition-all duration-500"
          />

          {/* Track outline */}
          <path
            d={result.d}
            stroke="oklch(0.35 0.02 260)"
            strokeWidth="8"
            strokeLinecap="round"
            strokeLinejoin="round"
            opacity="0.3"
          />

          {/* Start/Finish line indicator */}
          {result.firstPoint && (
            <>
              <circle
                cx={result.firstPoint.x}
                cy={result.firstPoint.y}
                r="6"
                fill="oklch(0.65 0.2 25)"
              />
              <circle
                cx={result.firstPoint.x}
                cy={result.firstPoint.y}
                r="3"
                fill="oklch(0.95 0.01 80)"
              />
            </>
          )}
        </svg>
      ) : (
        <div className="h-full w-full bg-background/5" />
      )}
    </div>
  );
}
