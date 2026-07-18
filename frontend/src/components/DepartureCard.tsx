import type { Departure } from "../types/timetable";
import { rowGridClasses, type BoardVariant } from "./boardLayout";

interface DepartureCardProps {
  departure: Departure;
  variant: BoardVariant;
}

export default function DepartureCard({
  departure,
  variant,
}: DepartureCardProps) {
  return (
    <div
      className={`${rowGridClasses[variant]} h-full min-h-0 border-b border-slate-700 py-2 last:border-b-0`}
    >
      <div className="min-w-0">
        <div className="truncate text-3xl font-bold text-white lg:text-4xl">
          {departure.routeName}
        </div>
      </div>

      <div className="min-w-0 text-center">
        <div className="truncate text-3xl font-bold text-white lg:text-4xl">
          {departure.destination}
        </div>
      </div>

      <div className="text-center font-mono text-4xl font-bold tabular-nums text-white lg:text-5xl">
        {departure.departureTime}
      </div>

      {variant === "bus" && (
        <div className="text-center text-2xl text-white lg:text-3xl">
          {departure.platform}
        </div>
      )}
    </div>
  );
}
