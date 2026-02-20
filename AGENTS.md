# [CRITICAL] PROJECT STANDARDS & ENFORCEMENT GUIDE

> **ATTENTION AGENT:** You are an automated contributor to the F1 Data platform. This file is your PRIMARY CONSTRAINT. Before every response, verify your plan against these rules.

---

## 1. PROJECT SCOPE & BOUNDARIES
* **ACTIVE:** `backend/` (Go), `fastf1_wrapper/` (Python).
* **STRICTLY IGNORE:** `frontend_remote/`. (Do not read, modify, or reference).
* **EXEMPT ZONES:** `/examples/` and `/scratch/`.
  - **Rule:** These folders are for experiments and prototypes. 
  - **Exemption:** Sections 2, 3, and 4 do NOT apply here. Prioritize speed/working code over standards in these paths.

---

## 2. ARCHITECTURE DECISIONS
### Service Roles
* **The Brain (`backend`):** Sole source of truth for orchestration, business logic, and **Data Formatting** (converting ms to human-readable strings).
* **The Adapter (`fastf1_wrapper`):** Fetch raw data -> Normalize types -> Pass to Go. **ZERO** complex domain calculations allowed here.

### Data Strategy
* **TIME TRANSPORT:** All durations/times between `fastf1_wrapper` and `backend` MUST be `int64` milliseconds (`ms`). 
* **DATA FORMATTING:** The backend MUST provide human-readable string versions of these `ms` values (e.g., `start_date` alongside `start_date_ms`) for consumer consumption.
* **COMPOSITION:** Use pointers with `omitempty` (Go) and `Optional[]` (Python) for clean, optional data structures.
* **SYNC REQUIREMENT:** Changes to Python dataclass outputs **MUST** be reflected in Go models immediately to prevent unmarshaling errors.
* **PERSISTENCE:** Use SQLite for dev/prod. SQLite driver: `modernc.org/sqlite` (CGO-free).
* **SCHEMA CONVENTIONS:** 
    - Use UUIDs (strings) for Primary Keys (except where composite keys make sense).
    - Enable WAL mode and Foreign Key constraints in SQLite.
    - Do NOT ever deviate from current schema without first explicitly confirming any changes.
    - NEVER change the schema SQL string without updating the relevant structs in `models`.

---

## 3. CODING STANDARDS (STRICT ENFORCEMENT)
*Applies to `backend/` and `fastf1_wrapper/` only.*

### Naming Conventions
* **FORBIDDEN:** Single-letter variables (`s`, `v`), cryptic abbreviations (`td`, `s_type`), and vague terms (`data`, `obj`).
* **REQUIRED:** Explicit names: `RaceWeekend`, `elapsed_duration`, `lap_time_ms`.
* **EXCEPTIONS:** `i` (loops), `ctx` (context), `err` (Go), `ms` (millisecond suffix), `t` and `tt` (testing), `x` and `y` (coordinates).

### Implementation
* **Logical Cohesion:** Functions MUST follow the Single Responsibility Principle. Do not create "mega-functions" that handle multiple logical steps (e.g., fetch, transform, and format in one block).
* **Decomposition:** Extract complex sub-steps into helper functions to keep the main flow readable. 
* **Organization:** - Move logic into specialized scripts (e.g., `extractors.py`, `converters.py`).
    - **Avoid "Junk Drawers":** Do not create a single `utils.py` with dozens of unrelated functions. 
    - **Grouping:** Group related helper functions into context-specific scripts (e.g., `lap_timing_utils.py`, `session_parsing.py`) to maintain a clean directory structure.* **Python Models:** Use `dataclasses` only. Always load sessions with `laps=True`.
* **Go Layout:** Follow standard `cmd/` and `internal/` structure.

### Comments
* **PURPOSE:** Only use comments to explain the 'why' behind non-obvious logic or specific domain requirements (e.g., coordinate system conversions, specific racing rules).
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
* **MANDATORY WORKFLOW:** After changes, you MUST verify:
  1. `go test ./...`
  2. `python3 -m unittest discover`

---

## 5. AGENT MAINTENANCE (THE SCRIBE RULE)
* **SELF-UPDATE:** You are responsible for this file. When new architectural decisions are made or patterns established, you MUST update Section 2 or 3 before concluding the session.
* **DEPENDENCIES:** NEVER install new libraries (pip/go get) without explaining the necessity and receiving explicit user confirmation.

---

## 6. SELF-CORRECTION PROTOCOL
If the user reports a "Standard Violation," you must:
1. Identify the specific rule in this file you breached.
2. Explain the error and refactor the code to 100% compliance immediately.
