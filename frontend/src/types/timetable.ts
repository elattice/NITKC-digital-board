export type DepartureKind = "bus" | "train";

export interface Departure {
  id: number;
  kind: DepartureKind;
  routeName: string;
  destination: string;
  departureTime: string;
  platform: string;
}

export interface UpcomingDeparturesResponse {
  bus: Departure[];
  train: Departure[];
  updatedAt: string;
}

export interface DeparturesResponse {
  departures: Departure[];
}

export type ImportCSVResult =
  | { ok: true; imported: number }
  | { ok: false; error: string; details: string[] };
