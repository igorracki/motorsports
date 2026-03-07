"use client";

import React from "react";
import Image from "next/image";
import { useState, useCallback } from "react";
import { GripVertical } from "lucide-react";
import { cn } from "@/lib/utils";
import type { DriverInfo } from "@/types/f1";

interface PredictionTableProps {
  drivers: DriverInfo[];
  onPredictionsChange: (predictions: DriverInfo[]) => void;
  initialPredictions?: DriverInfo[];
  onSave: (predictions: DriverInfo[]) => void;
}

export function PredictionTable({
  drivers,
  onPredictionsChange,
  initialPredictions,
  onSave,
}: PredictionTableProps) {
  const [predictions, setPredictions] = useState<DriverInfo[]>(
    initialPredictions || [...drivers]
  );

  // Notify parent of changes
  const updatePredictions = useCallback(
    (newPredictions: DriverInfo[]) => {
      setPredictions(newPredictions);
      onPredictionsChange(newPredictions);
    },
    [onPredictionsChange]
  );
  const [draggedIndex, setDraggedIndex] = useState<number | null>(null);
  const [dragOverIndex, setDragOverIndex] = useState<number | null>(null);

  const handleDragStart = useCallback(
    (e: React.DragEvent<HTMLTableRowElement>, index: number) => {
      setDraggedIndex(index);
      e.dataTransfer.effectAllowed = "move";
      e.dataTransfer.setData("text/plain", index.toString());
      // Add a slight delay to allow the drag image to be set
      setTimeout(() => {
        (e.target as HTMLTableRowElement).style.opacity = "0.5";
      }, 0);
    },
    []
  );

  const handleDragEnd = useCallback(
    (e: React.DragEvent<HTMLTableRowElement>) => {
      (e.target as HTMLTableRowElement).style.opacity = "1";
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

      const newPredictions = [...predictions];
      const [draggedItem] = newPredictions.splice(dragIndex, 1);
      newPredictions.splice(dropIndex, 0, draggedItem);
      updatePredictions(newPredictions);

      setDraggedIndex(null);
      setDragOverIndex(null);
    },
    [predictions, updatePredictions]
  );

  // Touch support for mobile
  const [touchStartY, setTouchStartY] = useState<number | null>(null);
  const [touchedIndex, setTouchedIndex] = useState<number | null>(null);

  const handleTouchStart = useCallback(
    (e: React.TouchEvent<HTMLTableRowElement>, index: number) => {
      setTouchStartY(e.touches[0].clientY);
      setTouchedIndex(index);
    },
    []
  );

  const handleTouchMove = useCallback(
    (e: React.TouchEvent<HTMLTableRowElement>) => {
      if (touchStartY === null || touchedIndex === null) return;

      const currentY = e.touches[0].clientY;
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
    },
    [touchStartY, touchedIndex]
  );

  const handleTouchEnd = useCallback(() => {
    if (touchedIndex !== null && dragOverIndex !== null && touchedIndex !== dragOverIndex) {
      const newPredictions = [...predictions];
      const [draggedItem] = newPredictions.splice(touchedIndex, 1);
      newPredictions.splice(dragOverIndex, 0, draggedItem);
      updatePredictions(newPredictions);
    }
    setTouchStartY(null);
    setTouchedIndex(null);
    setDragOverIndex(null);
  }, [touchedIndex, dragOverIndex, predictions, updatePredictions]);

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
              <th className="w-10 px-2 py-3"></th>
            </tr>
          </thead>
          <tbody>
            {predictions.map((driver, index) => (
              <tr
                key={driver.id}
                data-row-index={index}
                draggable
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
                  dragOverIndex === index &&
                    draggedIndex !== null &&
                    "bg-primary/10 border-primary/50",
                  dragOverIndex !== index && "hover:bg-secondary/30"
                )}
              >
                <td className="px-3 py-3 text-center">
                  <span
                    className={cn(
                      "inline-flex h-7 w-7 items-center justify-center rounded-full text-sm font-bold",
                      index === 0 &&
                        "bg-yellow-500/20 text-yellow-500",
                      index === 1 &&
                        "bg-gray-400/20 text-gray-400",
                      index === 2 &&
                        "bg-amber-600/20 text-amber-600",
                      index > 2 && "text-muted-foreground"
                    )}
                  >
                    {index + 1}
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
