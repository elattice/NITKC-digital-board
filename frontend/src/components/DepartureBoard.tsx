import type { Departure } from "../types/timetable";
import DepartureCard from "./DepartureCard";
import {
  accentBars,
  rowGridClasses,
  type BoardAccent,
  type BoardVariant,
} from "./boardLayout";

interface DepartureBoardProps {
  title: string;
  departures: Departure[];
  accent: BoardAccent;
  variant: BoardVariant;
}

const MAX_VISIBLE_DEPARTURES = 4;

const emptyMessages: Record<BoardVariant, string> = {
  bus: "本日のバス表示対象はありません",
  train: "本日の電車表示対象はありません",
};

function HeaderLabel({ ja, en }: { ja: string; en: string }) {
  return (
    <div className="text-center leading-tight">
      <div className="text-xl font-bold text-slate-200 lg:text-2xl">{ja}</div>
      <div className="text-xs text-slate-400 lg:text-sm">{en}</div>
    </div>
  );
}

export default function DepartureBoard({
  title,
  departures,
  accent,
  variant,
}: DepartureBoardProps) {
  const visibleDepartures = departures.slice(0, MAX_VISIBLE_DEPARTURES);

  return (
    <section className="flex min-h-0 flex-col overflow-hidden">
      <h2
        className={`${accentBars[accent]} shrink-0 px-6 py-1 text-2xl font-bold text-white lg:text-3xl`}
      >
        {title}
      </h2>

      {visibleDepartures.length === 0 ? (
        <div className="flex min-h-0 flex-1 items-center px-6 py-3">
          <p className="w-full rounded-lg border border-slate-700 bg-black px-8 py-6 text-center text-3xl font-bold text-slate-300 lg:text-4xl">
            {emptyMessages[variant]}
          </p>
        </div>
      ) : (
        <>
          <div
            className={`${rowGridClasses[variant]} shrink-0 border-b border-slate-600 py-1.5`}
          >
            <HeaderLabel ja="路線名" en="Route" />
            <HeaderLabel ja="行先" en="Destination" />
            <HeaderLabel ja="発車時刻" en="Dep. Time" />
            {variant === "bus" ? (
              <HeaderLabel ja="のりば" en="Platform" />
            ) : (
              <HeaderLabel ja="駅名" en="Station" />
            )}
          </div>

          <div
            className="grid min-h-0 flex-1 overflow-hidden"
            style={{
              gridTemplateRows: `repeat(${visibleDepartures.length}, minmax(0, 1fr))`,
            }}
          >
            {visibleDepartures.map((departure) => (
              <DepartureCard
                key={departure.id}
                departure={departure}
                variant={variant}
              />
            ))}
          </div>
        </>
      )}
    </section>
  );
}
