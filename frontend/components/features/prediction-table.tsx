"use client";

import React, { useMemo, useCallback, useState, useEffect } from "react";
import Image from "next/image";
import { GripVertical, Trophy } from "lucide-react";
import {
  DndContext,
  closestCenter,
  KeyboardSensor,
  PointerSensor,
  TouchSensor,
  useSensor,
  useSensors,
  DragEndEvent,
} from "@dnd-kit/core";
import {
  arrayMove,
  SortableContext,
  sortableKeyboardCoordinates,
  verticalListSortingStrategy,
  useSortable,
} from "@dnd-kit/sortable";
import { CSS } from "@dnd-kit/utilities";
import { restrictToVerticalAxis } from "@dnd-kit/modifiers";
import { cn } from "@/lib/utils";
import type { DriverInfo } from "@/types/f1";
import { LoadingSpinner } from "@/components/ui/LoadingSpinner";

interface PredictionTableProps {
  drivers: DriverInfo[];
  onPredictionsChange: (predictions: DriverInfo[]) => void;
  readOnly?: boolean;
  totalScore?: number;
}

interface SortableDriverRowProps {
  driver: DriverInfo;
  index: number;
  readOnly: boolean;
  onToggle: (index: number) => void;
}

function DriverRowContent({
  driver,
  index,
}: {
  driver: DriverInfo;
  index: number;
}) {
  return (
    <>
      <td className="px-3 py-3 text-center">
        <span
          className={cn(
            "inline-flex h-7 w-7 items-center justify-center rounded-full text-sm font-bold",
            !driver.isPredicted && "text-muted-foreground/30",
            driver.isPredicted &&
              index === 0 &&
              "bg-yellow-500/20 text-yellow-500",
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
        <span className="font-medium text-foreground">{driver.fullName}</span>
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
    </>
  );
}

function SortableDriverRow({
  driver,
  index,
  readOnly,
  onToggle,
}: SortableDriverRowProps) {
  const {
    attributes,
    listeners,
    setNodeRef,
    transform,
    transition,
    isDragging,
  } = useSortable({
    id: driver.id,
    disabled: readOnly,
  });

  const style = {
    transform: CSS.Translate.toString(transform),
    transition,
  };

  return (
    <tr
      ref={setNodeRef}
      style={style}
      {...attributes}
      {...listeners}
      onDoubleClick={() => onToggle(index)}
      className={cn(
        "cursor-grab border-b border-border/30 transition-all duration-150 active:cursor-grabbing",
        isDragging &&
          "opacity-50 scale-[1.02] shadow-xl bg-primary/20 z-10 relative",
        driver.isPredicted &&
          !driver.correct &&
          "bg-blue-500/5 border-blue-500/20",
        driver.correct && "bg-success/10 border-success/30",
        !isDragging && "hover:bg-secondary/30"
      )}
    >
      <DriverRowContent driver={driver} index={index} />
    </tr>
  );
}

export function PredictionTable({
  drivers,
  onPredictionsChange,
  readOnly = false,
  totalScore,
}: PredictionTableProps) {
  const [isMounted, setIsMounted] = useState(false);

  useEffect(() => {
    // Hydration guard: dnd-kit and many of our complex UI hooks use window sensors 
    // and generated IDs that won't match between server and client during hydration.
    // We render a simple skeleton first to ensure a stable hydration pass.
    // eslint-disable-next-line react-hooks/set-state-in-effect
    setIsMounted(true);
  }, []);

  // Notify parent of changes
  const updatePredictions = useCallback(
    (newPredictions: DriverInfo[]) => {
      if (readOnly) return;
      onPredictionsChange(newPredictions);
    },
    [onPredictionsChange, readOnly]
  );

  const sensors = useSensors(
    useSensor(PointerSensor, {
      activationConstraint: { distance: 8 },
    }),
    useSensor(TouchSensor, {
      activationConstraint: {
        delay: 250,
        tolerance: 5,
      },
    }),
    useSensor(KeyboardSensor, {
      coordinateGetter: sortableKeyboardCoordinates,
    })
  );

  const togglePrediction = useCallback(
    (index: number) => {
      if (readOnly) return;

      const newPredictions = [...drivers];
      newPredictions[index] = {
        ...newPredictions[index],
        isPredicted: !newPredictions[index].isPredicted,
      };

      updatePredictions(newPredictions);
    },
    [drivers, readOnly, updatePredictions]
  );

  const handleDragEnd = (event: DragEndEvent) => {
    const { active, over } = event;

    if (over && active.id !== over.id) {
      const oldIndex = drivers.findIndex((d) => d.id === active.id);
      const newIndex = drivers.findIndex((d) => d.id === over.id);

      const newDrivers = arrayMove([...drivers], oldIndex, newIndex);

      // Mark the moved item as predicted
      newDrivers[newIndex] = {
        ...newDrivers[newIndex],
        isPredicted: true,
      };

      updatePredictions(newDrivers);
    }
  };

  const driverIds = useMemo(() => drivers.map((d) => d.id), [drivers]);

  // If not mounted, render a stable, non-interactive skeleton 
  // to ensure server-side and client-side HTML matches perfectly.
  if (!isMounted) {
    return (
      <div className="flex h-96 items-center justify-center rounded-xl border border-border/50 bg-card/50">
        <LoadingSpinner label="Initializing prediction table..." />
      </div>
    );
  }

  return (
    <div className="overflow-hidden rounded-xl border border-border/50 bg-card">
      <div className="overflow-x-auto">
        <DndContext
          sensors={sensors}
          collisionDetection={closestCenter}
          onDragEnd={handleDragEnd}
          modifiers={[restrictToVerticalAxis]}
        >
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
                      <span className="text-[10px] text-muted-foreground font-semibold uppercase">
                        Total:
                      </span>
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
              <SortableContext
                items={driverIds}
                strategy={verticalListSortingStrategy}
              >
                {drivers.map((driver, index) => (
                  <SortableDriverRow
                    key={driver.id}
                    driver={driver}
                    index={index}
                    readOnly={readOnly}
                    onToggle={togglePrediction}
                  />
                ))}
              </SortableContext>
            </tbody>
          </table>
        </DndContext>
      </div>
    </div>
  );
}
