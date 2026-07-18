# Traffic Information Board

A system for displaying upcoming bus and train departure times on a large monitor.

- Display board (`/`): Shows upcoming departures in the style of an electronic information board.
- Admin page (`/admin`): Replaces the timetable data by importing a CSV file.

## Tech Stack

| Component | Technologies |
| --- | --- |
| Frontend | Vite / React / TypeScript / Tailwind CSS |
| Backend | Go / chi / SQLite (modernc.org/sqlite) |

Data flow: CSV → SQLite → Go API → React UI

## Project Structure

```text
digital-board/
├── backend/
│   ├── cmd/server/          # Entry point (router setup and startup)
│   ├── internal/db/         # SQLite connection, schema setup, and seed data
│   ├── internal/timetable/  # Types, database queries, API handlers, and CSV import
│   ├── internal/webui/      # Embedded React build with SPA fallback
│   └── data/timetable.db    # SQLite database (created on first run; not tracked by Git)
├── frontend/
│   ├── src/pages/           # BoardPage (display) / AdminPage (administration)
│   ├── src/components/      # Display-board UI components
│   ├── src/api/             # API client
│   └── vite.config.ts       # React build configuration
├── scripts/
│   └── build-release.sh     # Cross-build script for macOS and Ubuntu
├── dist/                    # Release binary output
└── docs/
    ├── build.md                          # Build instructions for macOS and Ubuntu
    ├── operation-manual.md               # Open-campus operation manual
    └── open-campus-holiday-timetable.csv # Example timetable CSV
```

## Requirements

- Go 1.26 or later
- Node.js 20 or later (including npm)

<a id="csv形式"></a>

## CSV Format

Use [docs/open-campus-holiday-timetable.csv](docs/open-campus-holiday-timetable.csv) as a reference when creating a timetable. Save the file as UTF-8 CSV. The first row must contain the following header in exactly this order:

```csv
kind,route_name,destination,departure_time,platform
```

Add one departure per row. For example:

```csv
kind,route_name,destination,departure_time,platform
bus,大楽毛線,釧路駅,07:40,高専前
train,根室本線,釧路方面,07:27,大楽毛駅
```

| Column | Required | Description |
| --- | --- | --- |
| `kind` | Yes | Transport type. Use only `bus` or `train`. |
| `route_name` | Yes | Route or line name shown on the display board. |
| `destination` | Yes | Destination shown on the display board. |
| `departure_time` | Yes | Departure time in zero-padded 24-hour `HH:MM` format, such as `07:40` or `15:25`. |
| `platform` | Yes | Bus stop or boarding location for buses; departure station name for trains. |

Important notes:

- Do not add, remove, rename, or reorder columns.
- Include every departure needed for the day in the same file.
- Importing a CSV replaces the entire existing timetable.
- Japanese text is supported when the file is saved as UTF-8. In Excel, select the UTF-8 CSV format when exporting.
- After importing the file from `/admin`, confirm that the number of imported rows and the displayed timetable are correct.
