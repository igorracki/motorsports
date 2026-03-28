"use client";

import React from "react";
import Image from "next/image";
import { useState, useCallback, useRef, useEffect } from "react";
import { GripVertical, Trophy } from "lucide-react";
import { cn } from "@/lib/utils";
import type { DriverInfo } from "@/types/f1";

interface PredictionTableProps {
  drivers: DriverInfo[];
  onPredictionsChange: (predictions: DriverInfo[]) => void;
  onSave: (predictions: DriverInfo[]) => void;
  readOnly?: boolean;
  totalScore?: number;
}

export function PredictionTable({
  drivers,
  onPredictionsChange,
  onSave,
  readOnly = false,
  totalScore,
}: PredictionTableProps) {
  // Notify parent of changes
  const updatePredictions = useCallback(
    (newPredictions: DriverInfo[]) => {
      if (readOnly) return;
      onPredictionsChange(newPredictions);
    },
    [onPredictionsChange, readOnly]
  );
  const [draggedIndex, setDraggedIndex] = useState<number | null>(null);
  const [dragOverIndex, setDragOverIndex] = useState<number | null>(null);

  // Touch support for mobile (Hold-to-drag)
  const [touchStartY, setTouchStartY] = useState<number | null>(null);
  const [touchedIndex, setTouchedIndex] = useState<number | null>(null);
  const [isLongPressActive, setIsLongPressActive] = useState(false);
  const touchTimeoutRef = useRef<NodeJS.Timeout | null>(null);
  const touchStartPositionRef = useRef<{ x: number, y: number } | null>(null);

  const clearTouchTimeout = useCallback(() => {
    if (touchTimeoutRef.current) {
      clearTimeout(touchTimeoutRef.current);
      touchTimeoutRef.current = null;
    }
  }, []);

  // Cleanup on unmount
  useEffect(() => {
    return () => clearTouchTimeout();
  }, [clearTouchTimeout]);

  const handleDoubleClick = useCallback((index: number) => {
    if (readOnly) return;

    const newPredictions = [...drivers];
    newPredictions[index] = {
      ...newPredictions[index],
      isPredicted: !newPredictions[index].isPredicted
    };
    
    updatePredictions(newPredictions);
  }, [drivers, readOnly, updatePredictions]);

  const handleDragStart = useCallback(
    (e: React.DragEvent<HTMLTableRowElement>, index: number) => {
      if (readOnly) {
        e.preventDefault();
        return;
      }
      setDraggedIndex(index);
      e.dataTransfer.effectAllowed = "move";
      e.dataTransfer.setData("text/plain", index.toString());
    },
    [readOnly]
  );

  const handleDragEnd = useCallback(
    (e: React.DragEvent<HTMLTableRowElement>) => {
      setDraggedIndex(null);
      setDragOverIndex(null);
    },
    []
  );

  const handleDragOver = useCallback(
    (e: React.DragEvent<HTMLTableRowElement>, index: number) => {
      e.preventDefault();
      e.dataTransfer.dropEffect = "move";
      if (draggedIndex !== null && draggedIndex !== index) {
        setDragOverIndex(index);
      }
    },
    [draggedIndex]
  );

  const handleDragLeave = useCallback(() => {
    setDragOverIndex(null);
  }, []);

  const handleDrop = useCallback(
    (e: React.DragEvent<HTMLTableRowElement>, dropIndex: number) => {
      e.preventDefault();
      const dragIndex = parseInt(e.dataTransfer.getData("text/plain"), 10);

      if (dragIndex === dropIndex) return;

      const newPredictions = [...drivers];
      const [draggedItem] = newPredictions.splice(dragIndex, 1);
      
      // Mark the moved item as predicted
      draggedItem.isPredicted = true;
      
      newPredictions.splice(dropIndex, 0, draggedItem);
      updatePredictions(newPredictions);

      setDraggedIndex(null);
      setDragOverIndex(null);
    },
    [drivers, updatePredictions]
  );

  // Touch support for double tap and hold-to-drag
  const [lastTap, setLastTap] = useState<{ time: number, index: number } | null>(null);

  const handleTouchStart = useCallback(
    (e: React.TouchEvent<HTMLTableRowElement>, index: number) => {
      if (readOnly) return;

      const now = Date.now();
      if (lastTap && now - lastTap.time < 300 && lastTap.index === index) {
        // Double tap detected
        handleDoubleClick(index);
        setLastTap(null);
        clearTouchTimeout();
        
        // IMPORTANT: Prevent default browser behavior for double-taps
        // This stops synthetic dblclick events from firing and causing a double-toggle
        if (e.cancelable) {
          e.preventDefault();
        }
        return;
      }
      setLastTap({ time: now, index });

      const touch = e.touches[0];
      setTouchStartY(touch.clientY);
      setTouchedIndex(index);
      touchStartPositionRef.current = { x: touch.clientX, y: touch.clientY };

      // Long press detection: 500ms
      clearTouchTimeout();
      touchTimeoutRef.current = setTimeout(() => {
        setIsLongPressActive(true);
        // Vibrate if supported
        if ("vibrate" in navigator) {
          navigator.vibrate(50);
        }
      }, 500);
    },
    [readOnly, lastTap, handleDoubleClick, clearTouchTimeout]
  );

  const handleTouchMove = useCallback(
    (e: React.TouchEvent<HTMLTableRowElement>) => {
      if (touchedIndex === null) return;

      const touch = e.touches[0];
      
      // Check for movement before long press threshold is reached
      if (!isLongPressActive && touchStartPositionRef.current) {
        const deltaX = Math.abs(touch.clientX - touchStartPositionRef.current.x);
        const deltaY = Math.abs(touch.clientY - touchStartPositionRef.current.y);
        
        // If they move significantly, assume they are scrolling and cancel long press
        if (deltaX > 10 || deltaY > 10) {
          clearTouchTimeout();
        }
      }

      if (isLongPressActive) {
        // Prevent default only during an active drag-and-drop
        if (e.cancelable) {
          e.preventDefault();
        }

        const currentY = touch.clientY;
        const rows = document.querySelectorAll("[data-row-index]");
        
        for (const row of rows) {
          const rect = row.getBoundingClientRect();
          if (currentY >= rect.top && currentY <= rect.bottom) {
            const newIndex = parseInt(
              row.getAttribute("data-row-index") || "0",
              10
            );
            if (newIndex !== touchedIndex) {
              setDragOverIndex(newIndex);
            }
            break;
          }
        }
      }
    },
    [touchedIndex, isLongPressActive, clearTouchTimeout]
  );

  const handleTouchEnd = useCallback(() => {
    clearTouchTimeout();
    
    if (isLongPressActive && touchedIndex !== null && dragOverIndex !== null && touchedIndex !== dragOverIndex) {
      const newPredictions = [...drivers];
      const [draggedItem] = newPredictions.splice(touchedIndex, 1);
      
      // Mark as predicted
      draggedItem.isPredicted = true;
      
      newPredictions.splice(dragOverIndex, 0, draggedItem);
      updatePredictions(newPredictions);
    }

    setTouchStartY(null);
    setTouchedIndex(null);
    setDragOverIndex(null);
    setIsLongPressActive(false);
    touchStartPositionRef.current = null;
  }, [clearTouchTimeout, isLongPressActive, touchedIndex, dragOverIndex, drivers, updatePredictions]);

  return (
    <div className="overflow-hidden rounded-xl border border-border/50 bg-card">
      <div className="overflow-x-auto">
        <table className="w-full">
          <thead>
            <tr className="border-b border-border/50 bg-secondary/50">
              <th className="w-12 px-3 py-3 text-center text-xs font-semibold uppercase tracking-wider text-muted-foreground">
                Pos
              </th>
              <th className="w-12 px-3 py-3 text-center text-xs font-semibold uppercase tracking-wider text-muted-foreground">
                Flag
              </th>
              <th className="px-4 py-3 text-left text-xs font-semibold uppercase tracking-wider text-muted-foreground">
                Driver
              </th>
              <th className="w-12 px-3 py-3 text-center text-xs font-semibold uppercase tracking-wider text-muted-foreground">
                #
              </th>
              <th className="px-4 py-3 text-left text-xs font-semibold uppercase tracking-wider text-muted-foreground">
                Team
              </th>
              <th className="w-20 px-3 py-3 text-right text-xs font-bold tracking-wider text-foreground">
                {totalScore !== undefined && totalScore !== null ? (
                  <div className="flex items-center justify-end gap-1.5 text-success">
                    <span className="text-[10px] text-muted-foreground font-semibold uppercase">Total:</span>
                    <span>{totalScore}</span>
                    <Trophy className="h-3 w-3 fill-success/20" />
                  </div>
                ) : (
                  "Points"
                )}
              </th>
              <th className="w-10 px-2 py-3"></th>
            </tr>
          </thead>
          <tbody>
            {drivers.map((driver, index) => (
              <tr
                key={driver.id}
                data-row-index={index}
                draggable={!readOnly}
                onDoubleClick={() => handleDoubleClick(index)}
                onDragStart={(e) => handleDragStart(e, index)}
                onDragEnd={handleDragEnd}
                onDragOver={(e) => handleDragOver(e, index)}
                onDragLeave={handleDragLeave}
                onDrop={(e) => handleDrop(e, index)}
                onTouchStart={(e) => handleTouchStart(e, index)}
                onTouchMove={handleTouchMove}
                onTouchEnd={handleTouchEnd}
                 className={cn(
                    "cursor-grab border-b border-border/30 transition-all duration-150 active:cursor-grabbing",
                    draggedIndex === index && "opacity-50",
                    touchedIndex === index && isLongPressActive && "opacity-50 scale-[1.02] shadow-xl bg-primary/20 z-10",
                    dragOverIndex === index &&
                      (draggedIndex !== null || isLongPressActive) &&
                      "bg-primary/10 border-primary/50",
                    dragOverIndex !== index && "hover:bg-secondary/30",
                    driver.isPredicted && !driver.correct && "bg-blue-500/5 border-blue-500/20",
                    driver.correct && "bg-success/10 border-success/30"
                  )}

               >

                <td className="px-3 py-3 text-center">
                  <span
                    className={cn(
                      "inline-flex h-7 w-7 items-center justify-center rounded-full text-sm font-bold",
                      !driver.isPredicted && "text-muted-foreground/30",
                      driver.isPredicted && index === 0 && "bg-yellow-500/20 text-yellow-500",
                      driver.isPredicted && index === 1 && "bg-gray-400/20 text-gray-400",
                      driver.isPredicted && index === 2 && "bg-amber-600/20 text-amber-600",
                      driver.isPredicted && index > 2 && "text-muted-foreground"
                    )}
                  >
                    {driver.isPredicted ? index + 1 : "-"}
                  </span>
                </td>
                <td className="px-3 py-3 text-center">
                  {driver.countryCode && (
                    <div className="mx-auto relative h-3 w-5 overflow-hidden rounded-sm ring-1 ring-border/50">
                      <Image
                        src={`https://flagcdn.com/w80/${driver.countryCode.toLowerCase()}.png`}
                        alt={`${driver.countryCode} flag`}
                        fill
                        className="object-cover"
                        unoptimized
                      />
                    </div>
                  )}
                </td>
                <td className="px-4 py-3">
                  <span className="font-medium text-foreground">
                    {driver.fullName}
                  </span>
                </td>
                <td className="px-3 py-3 text-center font-mono text-sm font-bold text-muted-foreground/80">
                  {driver.number}
                </td>
                 <td className="px-4 py-3 text-sm text-muted-foreground">
                   {driver.teamName}
                 </td>
                 <td className="px-3 py-3 text-right">
                   {driver.correct && driver.points > 0 && (
                     <div className="flex items-center justify-end gap-1 font-bold text-success animate-in fade-in zoom-in duration-500">
                       <span className="text-sm">+{driver.points}</span>
                       <Trophy className="h-3 w-3 fill-success/20" />
                     </div>
                   )}
                 </td>
                 <td className="px-2 py-3 text-right">

                  <GripVertical className="h-5 w-5 text-muted-foreground/30" />
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
}
