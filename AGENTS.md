# [CRITICAL] PROJECT STANDARDS & ENFORCEMENT GUIDE

> **ATTENTION AGENT:** You are an automated contributor to the F1 Data platform. This file is your PRIMARY CONSTRAINT. Before every response, verify your plan against these rules.

---

## 1. PROJECT SCOPE & BOUNDARIES
* **ACTIVE:** `backend/` (Go), `fastf1_wrapper/` (Python), `frontend/` (Next.js).
* **STRICTLY IGNORE:** `frontend_remote/`. (Do not read, modify, or reference).
* **EXEMPT ZONES:** `/examples/` and `/scratch/`.
  - **Rule:** These folders are for experiments and prototypes. 
  - **Exemption:** Sections 2, 3, and 4 do NOT apply here. Prioritize speed/working code over standards in these paths.

---

## 2. ARCHITECTURE DECISIONS
### Service Roles
* **The Brain (`backend`):** Sole source of truth for orchestration, business logic, and **Data Formatting** (converting ms to human-readable strings).
* **The Adapter (`fastf1_wrapper`):** Fetch raw data -> Normalize types -> Pass to Go. **ZERO** complex domain calculations allowed here.

### Security Standards (CRITICAL)
* **SQL INJECTION:** Always use parameterized queries (placeholders like `?`). Never use string concatenation for SQL.
* **XSS PROTECTION:** Sanitize all user-provided strings (e.g., `DisplayName`) before storage using `bluemonday`.
* **INPUT VALIDATION:** Implement strict allow-lists and bounds for all numeric parameters (e.g., years 1950-2100).
* **BODY LIMITS:** Enforce a maximum request body size (e.g., 1MB) to prevent DoS.
* **SECURE HEADERS:** Always enable secure HTTP headers (XSS Filter, Nosniff, etc.).
* **URL CONSTRUCTION:** Use proper URL builders (`url.JoinPath`, `url.URL`) and strictly validate path segments to prevent injection.

### Data Strategy
* **TIME TRANSPORT:** 
    - All durations/times between `fastf1_wrapper` and `backend` MUST be `int64` milliseconds (`ms`). 
    - `time_utc_ms` MUST be the true UTC epoch.
    - `utc_offset_ms` MUST be the track's local time offset from UTC at the time of the session.
* **GEOGRAPHIC DATA:**
    - `country_code` MUST be the ISO 3166-1 alpha-2 code (e.g., "BH", "US").
* **SESSION IDENTIFICATION:**
    - `session_code` MUST be the standard F1 abbreviation (e.g., "P1", "Q", "R").
* **CIRCUIT METRICS:**
    - Circuit data MUST include `max_speed_kmh` and elevation metrics (`max_altitude_m`, `min_altitude_m`) where telemetry is available.
* **DATA FORMATTING:** The backend MUST provide human-readable string versions of these `ms` values (e.g., `time_local` alongside `time_utc`) for consumer consumption. For weekend boundaries, both Local and UTC variants MUST be provided.
* **COMPOSITION:** Use pointers with `omitempty` (Go) and `Optional[]` (Python) for clean, optional data structures.
* **SYNC REQUIREMENT:** Changes to Python dataclass outputs **MUST** be reflected in Go models immediately to prevent unmarshaling errors.
* **PERSISTENCE:** Use SQLite for dev/prod. SQLite driver: `modernc.org/sqlite` (CGO-free).
* **SCHEMA CONVENTIONS:** 
    - Use UUIDs (strings) for Primary Keys (except where composite keys make sense).
    - Enable WAL mode and Foreign Key constraints in SQLite.
    - Do NOT ever deviate from current schema without first explicitly confirming any changes.
    - NEVER change the schema SQL string without updating the relevant structs in `models`.

### API Standards
* **COLLECTION CONSISTENCY:** Endpoints returning lists MUST return an empty array `[]` (never `null`) when no matches are found, accompanied by a `200 OK` status.
* **RESOURCE CONSISTENCY:** Endpoints fetching a single resource by identifier (e.g., `/users/:id`, `/circuits/:year/:round`) MUST return `404 Not Found` if the primary entity does not exist.
* **GO INITIALIZATION:** Always initialize slices in Go (e.g., `items := []models.Item{}`) before returning them in JSON responses to ensure `[]` marshaling.
* **DOCUMENTATION:** The OpenAPI specification (`backend/docs/openapi.yaml`) MUST be updated immediately after any change to API routes, request payloads, or response models.

### Frontend Architecture (Next.js)
* **SERVER FIRST:** Use Server Components for data fetching by default. Use Client Components only for interactivity.
* **RESTFUL ROUTING:** Follow the established path-segment pattern: `/race-weekend/[year]/[round]`. Do not use query parameters for primary resource identification.
* **DIRECTORY STRUCTURE:** 
    - `features/`: Complex, domain-specific dashboards and views.
    - `ui/`: Reusable, atomic UI primitives (Table, Skeleton, Badge).
    - `services/`: Dedicated API communication layer (`f1-api.ts`).
    - `hooks/`: Isolated client-side logic and global context (`SeasonContext`).
* **FRONTEND STATE:** Use `usePredictions` hook for synchronizing local UI state with backend persistence. Implement "Only Save Changes" logic by comparing current state against a fetched baseline.
* **TYPE MIRRORING:** Frontend interfaces in `types/f1.ts` MUST mirror Go models exactly. Use Zod transformations to map snake_case (Go) to camelCase (TS).
* **UI FEEDBACK:** Assignment-based tables (like predictions) MUST visually distinguish between "assigned" and "unassigned" states (e.g., using '-' placeholders and subtle color tints).

---

## 3. CODING STANDARDS (STRICT ENFORCEMENT)
*Applies to `backend/`, `fastf1_wrapper/`, and `frontend/`.*

### Naming Conventions
* **FORBIDDEN:** Single-letter variables (`s`, `v`), cryptic abbreviations (`td`, `s_type`), and vague terms (`data`, `obj`, `event`).
* **REQUIRED:** Explicit names: `RaceWeekend`, `raceWeekend`, `lapTimeMs`.
* **EXCEPTIONS:** `i` (loops), `ctx` (context), `err` (Go), `ms` (millisecond suffix), `t` and `tt` (testing), `x` and `y` (coordinates).

### Implementation
* **Logical Cohesion:** Functions MUST follow the Single Responsibility Principle. Do not create "mega-functions" that handle multiple logical steps (e.g., fetch, transform, and format in one block).
* **Decomposition:** Extract complex sub-steps into helper functions to keep the main flow readable. 
* **Organization:** - Move logic into specialized scripts (e.g., `extractors.py`, `converters.py`).
    - **Avoid "Junk Drawers":** Do not create a single `utils.py` with dozens of unrelated functions. 
    - **Grouping:** Group related helper functions into context-specific scripts (e.g., `lap_timing_utils.py`, `session_parsing.py`) to maintain a clean directory structure.
* **Python Models:** Use `dataclasses` only. Always load sessions with `laps=True`.
* **Go Layout:** Follow standard `cmd/` and `internal/` structure.
*   **Frontend UI:** 
    - **Optimization:** Always use the Next.js `<Image />` component for visual assets.
    - **UX:** Implement `loading.tsx` with custom `Skeleton` components for all dynamic routes.
    - **Consistency:** Use Tailwind classes only; avoid inline styles. Standardize colors via the project's design tokens.
    - **Dependencies:** Prefer built-in browser APIs (e.g., `Intl`, `fetch`, `crypto`) over external libraries (e.g., `date-fns`, `axios`) to minimize bundle size and leverage native performance.

### Comments
* **PURPOSE:** Only use comments to explain the 'why' behind non-obvious logic or specific domain requirements (e.g., coordinate system conversions, specific racing rules).
* **TESTING:** Always use "Given/When/Then" comments in behavioral tests.
* **FORBIDDEN:** Do not use comments to describe 'what' the code is doing if it is self-evident. Avoid 'step-by-step' numbering comments for standard flows.

### Observability & Logging
* **The "Log Sandwich" Pattern:** Every primary function must have an INFO log at the entry (input parameters) and the exit (result summary).
* **Traceability Requirements:**
    - **Context:** Always include identifiers (e.g., `driver_id`, `session_type`, `year`) in the log message.
    - **Counts:** When processing collections, log the count (e.g., "Processed 22 laps for Driver 44").
* **Error Handling:** - **Wrapper (Python):** Log the *original* exception and stack trace before returning a fallback value (`None`/`Empty`).
    - **Backend (Go):** Wrap errors with context: `fmt.Errorf("fetching weather for session %s: %w", sessionID, err)`.
* **FORBIDDEN:** Vague logs like "Done," "Error occurred," or "Success."

---

## 4. TESTING & VERIFICATION
* **POLICY:** Behavioral testing (Given/When/Then). Mock only at boundaries (e.g., HTTP calls).
* **COMMENTS:** Use "Given/When/Then" comments in all behavioral tests to clearly define the setup, action, and expected outcome.
* **MANDATORY WORKFLOW:** After changes, you MUST verify:
  1. `go test ./...`
  2. `python3 -m unittest discover`
  3. `cd frontend && npm run build && npx tsc --noEmit`

---

## 5. AGENT MAINTENANCE (THE SCRIBE RULE)
* **SELF-UPDATE:** You are responsible for this file. When new architectural decisions are made or patterns established, you MUST update Section 2 or 3 before concluding the session.
* **DEPENDENCIES:** NEVER install new libraries (pip/go get) without explaining the necessity and receiving explicit user confirmation.

---

## 6. SELF-CORRECTION PROTOCOL
If the user reports a "Standard Violation," you must:
1. Identify the specific rule in this file you breached.
2. Explain the error and refactor the code to 100% compliance immediately.
